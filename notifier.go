package goatee

import (
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

type connection struct {
	ws *websocket.Conn
	// buffered channel of outbound messages
	send chan []byte
}

type sockethub struct {
	// registered connections
	connections map[*connection]bool
	// inbound messages from connections
	broadcast chan []byte
	// register requests from connection
	register chan *connection
	// unregister request from connection
	unregister chan *connection
}

var h = sockethub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(1, message)
		if err != nil {
			log.Printf("Error in writer: ", err.Error())
			h.unregister <- c
			break
		}
	}
	c.ws.Close()
}

func LongPoll(w http.ResponseWriter, r *http.Request) {
	c := &connection{send: make(chan []byte, 256)}
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
		io.WriteString(w, string(msg))
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
		log.Printf("WsHandler error: ", err.Error())
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c

	defer func() { h.unregister <- c }()
	c.writer()
}

func (h *sockethub) Run() {
	for {
		select {
		case c := <-h.register:
			if DEBUG {
				log.Println("Connection created.")
			}
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
					if DEBUG {
						log.Printf("broadcasting: %s", string(m))
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
