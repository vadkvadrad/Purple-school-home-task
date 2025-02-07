package main

import (
	"fmt"
	"net/http"
	"verify-api/configs"
)

func main() {
	app := App()

	server := http.Server{
		Addr: ":8081",
		Handler: app,
	}

	fmt.Println("server started on port ", server.Addr)
	server.ListenAndServe()
}


func App() http.Handler {
	config := configs.Load()
	router := http.NewServeMux()

	// Handlers

	// Repository

	// Services

	return router
}