package goatee

import (
	"github.com/garyburd/redigo/redis"
	"github.com/johnernaut/goatee/config"
	"log"
	"sync"
	"time"
)

type Message struct {
	Type    string
	Channel string
	Data    string
}

type RedisClient struct {
	conn redis.Conn
	redis.PubSubConn
	sync.Mutex
}

type Client interface {
	Publish(channel, message string)
	Receive() (message Message)
}

func NewRedisClient(host, sub string) *RedisClient, error {
	conn, err := redis.Dial("tcp", host)

	if err != nil {
		log.Printf("Error dialing redis pubsub: %s", err)
		return nil, err
	}

	pubsub, _ := redis.Dial("tcp", host)

	client := RedisClient{conn, redis.PubSubConn{pubsub}, sync.Mutex{}}

	log.Println("Subscribed to Redis on: ", host)

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			client.Lock()
			client.conn.Flush()
			client.Unlock()
		}
	}()

	go client.PubsubHub()

	client.Subscribe(sub)

	return &client, nil
}

func (client *RedisClient) Receive() Message {
	switch message := client.PubSubConn.Receive().(type) {
	case redis.Message:
		return Message{"message", message.Channel, string(message.Data)}
	case redis.Subscription:
		return Message{message.Kind, message.Channel, string(message.Count)}
	}

	return Message{}
}

func (client *RedisClient) Publish(channel, message string) {
	client.Lock()
	client.conn.Send("PUBLISH", channel, message)
	client.Unlock()
}

func (client *RedisClient) PubsubHub() {
	for {
		message := client.Receive()
		if message.Type == "message" {
			H.Broadcast <- []byte(message.Data)
			if config.DEBUG {
				log.Printf("Received: %s", message.Data)
			}
		}
	}
}
