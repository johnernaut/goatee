package goatee

func CreateServer() error {
	conf := LoadConfig("config/")
	client, err := NewRedisClient(conf.Redis.Host)
	if err != nil {
		return err
	}

	defer client.Close()

	return NotificationHub(conf.Web.Host)
}
