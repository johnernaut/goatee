package main

import (
	"github.com/johnernaut/goatee"
	"log"
)

func main() {
	err := goatee.CreateServer([]string{"chan1", "chan2"})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
}
