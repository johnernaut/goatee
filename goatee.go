package goatee

func CreateServer(redisub string) error {
	client, err := NewRedisClient(Config.Redis.Host, redisub)
	if err != nil {
		return err
	}

	defer client.Close()

	// client.Publish("supdood", "a message from golang")
	// client.Subscribe("supdood")

	return NotificationHub(Config.Web.Host)
}
