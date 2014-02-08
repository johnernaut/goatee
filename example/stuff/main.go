package main

import (
    "github.com/johnernaut/goatee"
    "log"
)

func main() {
    err := goatee.CreateServer("achannel")

    if err != nil {
        log.Fatal("Error: ", err.Error())
    }
}
