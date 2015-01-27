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
    var file[]byte
    var err error
    var paths = []string{os.Getenv("HOME") + "/.config/goatee", "/etc/goatee"}

    // If path is defined, prepend it to paths
    if (len(path) > 0) {
        paths = append([]string{path}, paths...)
    }

    // Try to find a config file to use
    found := false
    for _, path := range(paths) {
        file, err = ioutil.ReadFile(path + string(os.PathSeparator) + getEnv() + ".json")
        if err == nil {
            log.Printf("Reading configuration from: %s", path)
            found = true
            break
        }
    }

    if !found {
        log.Fatalf("Error reading config file.")
    }

	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err.Error())
	}

	return Config
}
