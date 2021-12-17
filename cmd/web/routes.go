package main

import (
	"net/http"
	"ws/internal/handlers"

	"github.com/bmizerany/pat"
)

func routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))
	fileServer := http.FileServer(http.Dir("../internal/static/"))
	mux.Get("/static/", http.StripPrefix("../internal/static/", fileServer))
	return mux
}
