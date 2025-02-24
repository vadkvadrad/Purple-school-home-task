package main

import (
	"fmt"
	"net/http"
	"verify-api/configs"
	"verify-api/internal/user"
	"verify-api/internal/verify"
	"verify-api/pkg/db"
	"verify-api/pkg/middleware"
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
	config, err := configs.Load()
	if err != nil {
		panic(err)
	}
	db := db.NewDb(config)
	router := http.NewServeMux()

	// Repository
	userRepository := user.NewUserRepository(db)

	// Services
	verifyService := verify.NewVerifyService(verify.VerifyServiceDeps{
		UserRepository: userRepository,
		Config: config,
	})

	// Handlers
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		VerifyService: verifyService,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.Logging,
	)

	return stack(router)
}