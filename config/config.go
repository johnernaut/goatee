package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func getEnv() string {
	if os.Getenv("GO_ENV") == "" {
		return "development"
	} else {
		return os.Getenv("GO_ENV")
	}
}

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

var Config = new(configuration)

func init() {
	file, err := ioutil.ReadFile("config/" + getEnv() + ".json")

	if err != nil {
		log.Fatalf("Error parsing config: %s", err.Error())
	}

	err = json.Unmarshal(file, &Config)

	if err != nil {
		log.Fatalf("Error parsing json: %s", err.Error())
	}
}
