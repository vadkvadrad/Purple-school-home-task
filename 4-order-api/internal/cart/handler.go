package cart

import (
	"net/http"
	"order-api/configs"
)

type CartHandler struct {
	Config *configs.Config
	CartService *CartService
}

type CartHandlerDeps struct {
	Config *configs.Config
	CartService *CartService
}

func NewCartHandler(router *http.ServeMux, deps CartHandlerDeps) {
	handler := &CartHandler{
		Config: deps.Config,
		CartService: deps.CartService,
	}

	// Создание нового заказа
	router.HandleFunc("POST /order", handler.Create())

	// Получение заказа по ID
	router.HandleFunc("GET /order/{id}", handler.GetById())

	// Получение заказа по пользователю
	router.HandleFunc("GET /my-orders", handler.GetAll())
}


func (handler *CartHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}

func(handler *CartHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}

func(handler *CartHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}