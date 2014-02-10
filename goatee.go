package goatee

func CreateServer(redisub []string) error {
	conf := LoadConfig("config/")
	client, err := NewRedisClient(conf.Redis.Host, redisub)
	if err != nil {
		return err
	}

	defer client.Close()

	return NotificationHub(conf.Web.Host)
}
