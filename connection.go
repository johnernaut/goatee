package goatee

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
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
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if c.client.Channel == message.Channel {
				err := c.ws.WriteJSON(message.Payload)
				if err != nil {
					log.Printf("Error in writer: %s", err.Error())
					h.unregister <- c
					break
				}
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *connection) reader() {
	h.rclient.Lock()

	defer func() {
		h.unregister <- c
		h.rclient.Unlock()
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

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
		case "message":
			d, err := json.Marshal(wclient)
			if err != nil {
				log.Println("Error marsahling json for publish: ", err)
			}

			_, err = h.rconn.Do("PUBLISH", wclient.Channel, d)
			if err != nil {
				log.Println("Error publishing message: ", err)
			}
		}
	}
}
