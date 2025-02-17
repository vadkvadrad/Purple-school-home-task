package main

import (
	"fmt"
	"net/http"
	"order-api/internal/product"
	"order-api/pkg/middleware"
	"order-api/configs"
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
	conf, err := configs.Load()
	if err != nil {
		panic(err)
	}
	router := http.NewServeMux()

	// Handlers
	product.NewProductHandler(router, &product.ProductHandlerDeps{
		Config: conf,
	})

	// Middleware
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}
