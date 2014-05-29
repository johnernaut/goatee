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
	send   chan *Data
	client WSClient
}

func (c *connection) writer() {
	for message := range c.send {
		if c.client.Channel == message.Channel {
			err := c.ws.WriteMessage(1, []byte(message.Payload))
			if err != nil {
				log.Printf("Error in writer: %s", err.Error())
				h.unregister <- c
				break
			}
		}
	}
	c.ws.Close()
}

func (c *connection) reader() {
	for {
		var wclient WSClient
		err := c.ws.ReadJSON(&wclient)
		if err != nil {
			break
		}

		if DEBUG {
			log.Println("client type is:", wclient)
		}

		c.client = wclient

		switch wclient.Action {
		case "bind":
			h.rclient.Subscribe(wclient.Channel)
		case "unbind":
			h.rclient.Unsubscribe(wclient.Channel)
		}
	}

	c.ws.Close()
}
