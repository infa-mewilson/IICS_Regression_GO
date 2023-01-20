package main

import (
	"Golangcode/apmservice"
	"log"
	"net/http"
)

func main() {
	apmservice.Init()
	log.Println("Starting server on port 4048")
	log.Println(http.ListenAndServe(":4048", nil))
}
