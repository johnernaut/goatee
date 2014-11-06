package main

import (
	"log"
	"net/http"

	"github.com/johnernaut/goatee"
)

func main() {
	server := goatee.CreateServer()
	server.RegisterAuthFunc(Authenticate)
	server.StartServer()
}

func Authenticate(req *http.Request) bool {
	vals := req.URL.Query()

	if vals.Get("api_key") == "ABC123" {
		log.Println(vals.Get("api_key"))

		return true
	}

	return false
}
