package main

import (
	"fmt"
	"net/http"
	"verify-api/configs"
	"verify-api/internal/user"
	"verify-api/internal/verify"
	"verify-api/pkg/db"
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
	db := db.NewDb(config)
	router := http.NewServeMux()

	// Repository
	userRepository := user.NewUserRepository(db)

	// Services
	verifyService := verify.NewVerifyService(userRepository)

	// Handlers
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		VerifyService: verifyService,
	})

	return router
}