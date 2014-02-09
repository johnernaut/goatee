package goatee

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type configuration struct {
	Redis Redis
	Web   Web
}

type Redis struct {
	Host string
}

type Web struct {
	Host string
}

var (
	DEBUG  = false
	Config = new(configuration)
)

func getEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" || env == "development" {
		DEBUG = true
		return "development"
	}
	return env
}

func LoadConfig(path string) *configuration {
	file, err := ioutil.ReadFile(path + getEnv() + ".json")
	if err != nil {
		log.Fatalf("Error parsing config: %s", err.Error())
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatalf("Error parsing json: %s", err.Error())
	}

	return Config
}
