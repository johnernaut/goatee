package goatee

import (
	"github.com/gorilla/websocket"
	"github.com/johnernaut/goatee/config"
	"log"
	"net/http"
)

type connection struct {
	ws *websocket.Conn
	// buffered channel of outbound messages
	send chan []byte
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteJSON(string(message))
		if err != nil {
			log.Printf("Error in writer: ", err.Error())
		}
	}

	c.ws.Close()
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Printf("WsHandler error: ", err.Error())
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	H.register <- c
	//defer func() { H.unregister <- c }()
	go c.writer()
}

type sockethub struct {
	// registered connection
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

func (h *sockethub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.Broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
					if config.DEBUG {
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

	http.HandleFunc("/", WsHandler)

	log.Println("Starting websocket server on: ", host)

	if err := http.ListenAndServe(host, nil); err != nil {
		return error
	}

	return nil
}
