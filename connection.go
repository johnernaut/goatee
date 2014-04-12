package goatee

import (
	"github.com/gorilla/websocket"
	"log"
)

type connection struct {
	sid  string
	ws   *websocket.Conn
	send chan []byte
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(1, message)
		if err != nil {
			log.Printf("Error in writer: %s", err.Error())
			h.unregister <- c
			break
		}
	}
	c.ws.Close()
}
