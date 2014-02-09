package goatee

import (
	"testing"
)

type testData struct {
	host    string
	subchan string
}

func setup(t *testing.T) {
	setupRedisConnection(t)
	setupConfig(t)
}

func setupRedisConnection(t *testing.T) {
	data := testData{host: ":6379", subchan: "c1"}
	_, err := NewRedisClient(data.host, data.subchan)
	if err != nil {
		t.Errorf("Error creating Redis client: %s", err)
	}
}

func setupConfig(t *testing.T) {
	LoadConfig("fixture/")
}
