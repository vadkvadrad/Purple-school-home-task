package cart

import (
	"net/http"
	"order-api/configs"
	"order-api/internal/product"
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
	router.Handle("POST /order", middleware.IsAuthed(handler.Create(), handler.Config))

	// Получение заказа по ID
	router.Handle("GET /order/{id}", middleware.IsAuthed(handler.GetCartID(), handler.Config))

	// Получение заказа по пользователю
	router.Handle("GET /my-orders", middleware.IsAuthed(handler.GetByPhone(), handler.Config))
}


func (handler *CartHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, ErrNotAuthorized, http.StatusUnauthorized)
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

func(handler *CartHandler) GetByPhone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}

func(handler *CartHandler) GetCartID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, middleware.ErrUnauthorized, http.StatusUnauthorized)
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


		products := make([]product.Product, len(cart.Products))
		for i, productId := range cart.Products {
			prodId, err := strconv.Atoi(productId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			product, err := handler.CartService.ProductRepository.FindById(uint(prodId))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products[i] = *product
		}

		res.Json(w, products, http.StatusOK)
	}
}