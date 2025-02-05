package main

import (
	"fmt"
	"net/http"
	"rand-api/internal/random"
)

func main() {
	app := App()

	// creating server
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	// service start
	fmt.Println("service started on port", server.Addr)
	server.ListenAndServe()
}

func App() http.Handler {
	router := http.NewServeMux()

	// Handlers
	random.NewRandomHandler(router)

	return router
}