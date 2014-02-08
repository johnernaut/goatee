package goatee

import (
	"github.com/johnernaut/goatee/config"
	"github.com/johnernaut/goatee/goatee"
)

func CreateServer(redisub string) error {
	client, err := goatee.NewRedisClient(config.Config.Redis.Host, redisub)

	defer client.Close()

	// client.Publish("supdood", "a message from golang")
	// client.Subscribe("supdood")

	err = goatee.NotificationHub(config.Config.Web.Host)

	return err
}
