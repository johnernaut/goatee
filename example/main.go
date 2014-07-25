package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/johnernaut/goatee"
)

func main() {
	go httpServer()

	err := goatee.CreateServer()

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("chan1.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func httpServer() {
	log.Print("starting http server...")
	http.HandleFunc("/test", testHandler)
	log.Fatal(http.ListenAndServe(":1236", nil))
}
