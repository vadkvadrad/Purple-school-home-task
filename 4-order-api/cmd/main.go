package main

import (
	"fmt"
	"net/http"
	"order-api/configs"
	"order-api/internal/auth"
	"order-api/internal/user"
	"order-api/pkg/db"
	"order-api/pkg/middleware"
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
	db := db.NewDb(conf)

	// Repositories
	userRepository := user.NewUserRepository(db, conf)

	// Services
	authService := auth.NewAuthService(conf, userRepository)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
		AuthService: authService,
	})

	// Middleware
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}
