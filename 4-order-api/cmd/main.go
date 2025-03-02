package main

import (
	"fmt"
	"net/http"
	"order-api/configs"
	"order-api/internal/auth"
	"order-api/internal/cart/order"
	"order-api/internal/product/prod"
	"order-api/internal/user"
	"order-api/pkg/db"
	"order-api/pkg/middleware"
	"order-api/pkg/sender"
)

func main() {
	app := App()

	server := http.Server{
		Addr:    ":8081",
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
	sender, err := sender.Load(conf)
	if err != nil {
		panic(err)
	}

	// Repositories
	userRepository := user.NewUserRepository(db)
	cartRepository := order.NewCartRepository(db)
	productRepository := prod.NewProductRepository(db)

	// Services
	authService := auth.NewAuthService(auth.AuthServiceDeps{
		UserRepository: userRepository,
		Config:         conf,
		Sender:         sender,
	})
	cartService := order.NewCartService(order.CartServiceDeps{
		CartRepository:    cartRepository,
		ProductRepository: productRepository,
	})
	productService := prod.NewProductService(prod.ProductServiceDeps{
		ProductRepository: productRepository,
		CartRepository: cartRepository,
		UserRepository: userRepository,
		Sender:            sender,
	})

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	order.NewCartHandler(router, order.CartHandlerDeps{
		Config:         conf,
		CartService:    cartService,
		ProductService: productService,
	})
	prod.NewProductHandler(router, prod.ProductHandlerDeps{
		Config:         conf,
		ProductService: productService,
	})

	// Middleware
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}
