package main

import (
	"fmt"
	"net/http"
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
	router := http.NewServeMux()

	// Repository

	// Services

	return router
}