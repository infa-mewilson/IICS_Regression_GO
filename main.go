package main

import (
	"Golangcode/apmservice"
	"log"
	"net/http"
)

func main() {
	apmservice.Init()
	log.Println("Starting server on port 6666")
	log.Println(http.ListenAndServe(":6666", nil))
}
