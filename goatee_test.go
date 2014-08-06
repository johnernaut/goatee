package goatee

import (
	"testing"
)

type testData struct {
	host    string
	subchan []string
}

func setup(t *testing.T) *RedisClient {
	client := setupRedisConnection(t)
	setupConfig(t)

	return client
}

func setupRedisConnection(t *testing.T) *RedisClient {
	data := testData{host: ":6379"}
	client, err := NewRedisClient(data.host)
	if err != nil {
		t.Errorf("Error creating Redis client: %s", err)
	}

	return client
}

func setupConfig(t *testing.T) {
	LoadConfig("fixture/")
}
