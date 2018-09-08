package goatee

import (
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

// type AuthFunc func(req *http.Request) bool

type sockethub struct {
	// registered connections
	connections map[*connection]bool

	// inbound messages from connections
	broadcast chan *Data

	// register requests from connection
	register chan *connection

	// unregister request from connection
	unregister chan *connection

	// copy of the redis client
	rclient *RedisClient

	// copy of the redis connection
	rconn redis.Conn

	Auth func(req *http.Request) bool
}

var h = sockethub{
	broadcast:   make(chan *Data),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *sockethub) WsHandler(w http.ResponseWriter, r *http.Request) {
	var authenticated bool

	if h.Auth != nil {
		authenticated = h.Auth(r)
	}

	if authenticated {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
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
	} else {
		http.Error(w, "Invalid API key", 401)
	}
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

func (h *sockethub) RegisterAuthFunc(AuthFunc func(req *http.Request) bool) {
	aut := AuthFunc
	h.Auth = aut
}

func (h *sockethub) StartServer() {
	conf := LoadConfig("config")
	client, err := NewRedisClient(conf.Redis.Host)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	go h.Run()

	http.HandleFunc("/", h.WsHandler)
	log.Println("Starting server on: ", conf.Web.Host)

	err = http.ListenAndServe(conf.Web.Host, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateServer() sockethub {
	return h
}
