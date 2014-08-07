package goatee

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

type sockethub struct {
	// registered connections
	connections map[*connection]bool
	// inbound messages from connections
	broadcast chan *Data
	// register requests from connection
	register chan *connection
	// unregister request from connection
	unregister chan *connection
	rclient    *RedisClient
	rconn      redis.Conn
}

var h = sockethub{
	broadcast:   make(chan *Data),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func LongPoll(w http.ResponseWriter, r *http.Request) {
	c := &connection{send: make(chan *Data)}
	h.register <- c

	w.Header().Set("Access-Control-Allow-Origin", "*")

	cn, _ := w.(http.CloseNotifier)

	select {
	case <-time.After(30e9):
		io.WriteString(w, "Timeout!\n")
		h.unregister <- c
	case <-cn.CloseNotify():
		h.unregister <- c
	case msg := <-c.send:
		io.WriteString(w, msg.Payload)
		h.unregister <- c
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if _, ok := err.(websocket.HandshakeError); ok {
		// http.Error(w, "Not a websocket handshake", 400)
		LongPoll(w, r)
		return
	} else if err != nil {
		log.Printf("WsHandler error: %s", err.Error())
		return
	}

	c := &connection{send: make(chan *Data), ws: ws}
	h.register <- c

	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}

func (h *sockethub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
					if DEBUG {
						log.Printf("broadcasting: %s", m.Payload)
					}
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}

func NotificationHub(host string) error {
	go h.Run()
	http.HandleFunc("/", WsHandler)
	log.Println("Starting server on: ", host)
	return http.ListenAndServe(host, nil)
}
