package goatee

import (
	"github.com/gorilla/websocket"
	"log"
)

type WSClient struct {
	Channel string `json:"channel"`
	Action  string `json:"action"`
	Date    string `json:"date"`
	Payload string `json:"payload"`
	Token   string `json:"token"`
}

type connection struct {
	sid    string
	ws     *websocket.Conn
	send   chan []byte
	client WSClient
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

func (c *connection) reader() {
	for {
		var client WSClient
		err := c.ws.ReadJSON(&client)
		if err != nil {
			log.Println("[ERROR] we done messed up. ", err)
			break
		}

		if DEBUG {
			log.Println("client type is:", client)
		}

		c.client = client
	}
}
