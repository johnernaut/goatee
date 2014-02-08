package goatee

import (
	"github.com/johnernaut/goatee/config"
	"github.com/johnernaut/goatee/goatee"
)

func CreateServer(redisub string) {
	client := goatee.NewRedisClient(config.Config.Redis.Host, redisub)

	defer client.Close()

	// client.Publish("supdood", "a message from golang")
	// client.Subscribe("supdood")

	goatee.NewNotificationHub(config.Config.Web.Host)
}
