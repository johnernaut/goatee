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
	Broadcast chan []byte
	// register requests from connection
	register chan *connection
	// unregister request from connection
	unregister chan *connection
}

var H = sockethub{
	Broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(1, message)
		if err != nil {
			log.Printf("Error in writer: ", err.Error())
			H.unregister <- c
			break
		}
	}
	c.ws.Close()
}

func LongPoll(w http.ResponseWriter, r *http.Request) {
	c := &connection{send: make(chan []byte, 256)}
	H.register <- c

	w.Header().Set("Access-Control-Allow-Origin", "*")

	cn, _ := w.(http.CloseNotifier)

	select {
	case <-time.After(30e9):
		io.WriteString(w, "Timeout!\n")
		H.unregister <- c
	case <-cn.CloseNotify():
		H.unregister <- c
	case msg := <-c.send:
		io.WriteString(w, string(msg))
		H.unregister <- c
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
	H.register <- c

	defer func() { H.unregister <- c }()
	c.writer()
}

func (h *sockethub) Run() {
	for {
		select {
		case c := <-h.register:
			log.Println("Connection created.")
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.Broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
					if DEBUG {
						log.Printf("Broadcasting: %s", string(m))
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
	go H.Run()
	http.HandleFunc("/", LongPoll)
	log.Println("Starting server on: ", host)
	return http.ListenAndServe(host, nil)
}
