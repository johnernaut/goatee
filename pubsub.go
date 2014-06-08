package goatee

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"sync"
	"time"
)

type Message struct {
	Type    string
	Channel string
	Data    []byte
}

type Data struct {
	Channel   string `json:"channel"`
	Payload   string `json:"payload"`
	CreatedAt string `json:"created_at"`
}

type RedisClient struct {
	conn redis.Conn
	redis.PubSubConn
	sync.Mutex
}

type Client interface {
	Receive() (message Message)
}

func NewRedisClient(host string, sub []string) (*RedisClient, error) {
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		log.Printf("Error dialing redis pubsub: %s", err)
		return nil, err
	}

	pubsub, _ := redis.Dial("tcp", host)
	client := RedisClient{conn, redis.PubSubConn{pubsub}, sync.Mutex{}}

	if DEBUG {
		log.Println("Subscribed to Redis on: ", host)
	}

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			client.Lock()
			client.conn.Flush()
			client.Unlock()
		}
	}()

	go client.PubsubHub()

	// subscribe to default channels
	for _, k := range sub {
		client.PSubscribe(k)
	}

	h.rclient = &client
	h.rconn = conn

	return &client, nil
}

func (client *RedisClient) Receive() Message {
	switch message := client.PubSubConn.Receive().(type) {
	case redis.Message:
		return Message{"message", message.Channel, message.Data}
	case redis.Subscription:
		return Message{message.Kind, message.Channel, []byte(strconv.Itoa(message.Count))}
	}
	return Message{}
}

func (client *RedisClient) PubsubHub() {
	data := Data{}
	for {
		message := client.Receive()
		if message.Type == "message" {
			log.Println(string(message.Data))
			err := json.Unmarshal(message.Data, &data)
			if err != nil {
				log.Println("Error parsing payload JSON: ", err)
			}

			h.broadcast <- &data
			if DEBUG {
				log.Printf("Received: %s", message)
			}
		}
	}
}
