package main

import (
	"log"
	"net/http"
	"ws/internal/handlers"
)

func main() {

	mux := routes()
	log.Println("Starting Channel Listenr")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on Port 8080")

	_ = http.ListenAndServe(":8080", mux)
}