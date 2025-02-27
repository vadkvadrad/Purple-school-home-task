package cart

import (
	"net/http"
	"order-api/configs"
	"order-api/internal/product"
	"order-api/pkg/er"
	"order-api/pkg/middleware"
	"order-api/pkg/req"
	"order-api/pkg/res"
	"strconv"
	"time"

	"gorm.io/datatypes"
)

type CartHandler struct {
	Config *configs.Config
	CartService *CartService
	ProductRepository *product.ProductRepository
}

type CartHandlerDeps struct {
	Config *configs.Config
	CartService *CartService
	ProductRepository *product.ProductRepository
}

func NewCartHandler(router *http.ServeMux, deps CartHandlerDeps) {
	handler := &CartHandler{
		Config: deps.Config,
		CartService: deps.CartService,
		ProductRepository: deps.ProductRepository,
	}

	// Создание нового заказа
	router.Handle("POST /order", middleware.IsAuthed(handler.Create(), handler.Config))

	// Получение заказа по ID
	router.Handle("GET /order/{id}", middleware.IsAuthed(handler.GetByID(), handler.Config))

	// Получение заказа по пользователю
	router.Handle("GET /my-orders", middleware.IsAuthed(handler.GetByPhone(), handler.Config))

	// Обновление заказа
	router.Handle("PATCH /order/{id}", middleware.IsAuthed(handler.Update(), handler.Config))

	// Удаление заказа
	router.Handle("DELETE /order/{id}", middleware.IsAuthed(handler.Delete(), handler.Config))
}


func (handler *CartHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, er.ErrNotAuthorized, http.StatusUnauthorized)
			return
		}

		body, err := req.HandleBody[OrderRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cart, err := handler.CartService.Create(&Cart{
			Phone: phone,
			Products: body.Products,
			Date: datatypes.Date(time.Now()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, cart, http.StatusCreated)
	}
}

func(handler *CartHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, er.ErrNotAuthorized, http.StatusUnauthorized)
			return
		}

		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cart, err := handler.CartService.GetByIDAndPhone(id, phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		products, err := handler.CartService.ProductRepository.GetByIDs(cart.Products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, products, http.StatusOK)
	}
}

func(handler *CartHandler) GetByPhone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}

func(handler *CartHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}

func(handler *CartHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}