package main

import (
	"fmt"
	"net/http"
//	"order-api/configs"
)

func main() {
	app := App()

	server := http.Server{
		Addr: ":8081",
		Handler: app,
	}

	fmt.Println("server started on port", server.Addr)
	server.ListenAndServe()
}


func App() http.Handler {
	//conf, err := configs.Load()
	// if err != nil {
	// 	panic(err)
	// }

	router := http.NewServeMux()

	return router
}
