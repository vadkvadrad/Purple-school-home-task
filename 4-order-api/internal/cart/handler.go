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
	"gorm.io/gorm"
)

type CartHandler struct {
	Config         *configs.Config
	CartService    *CartService
	ProductService *product.ProductService
}

type CartHandlerDeps struct {
	Config         *configs.Config
	CartService    *CartService
	ProductService *product.ProductService
}

func NewCartHandler(router *http.ServeMux, deps CartHandlerDeps) {
	handler := &CartHandler{
		Config:         deps.Config,
		CartService:    deps.CartService,
		ProductService: deps.ProductService,
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
			Phone:    phone,
			Products: body.Products,
			Date:     datatypes.Date(time.Now()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, cart, http.StatusCreated)
	}
}

func (handler *CartHandler) GetByID() http.HandlerFunc {
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

		cart, err := handler.CartService.GetByID(id, phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		products, err := handler.ProductService.GetByIDs(cart.Products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, products, http.StatusOK)
	}
}

func (handler *CartHandler) GetByPhone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, er.ErrNotAuthorized, http.StatusUnauthorized)
		}

		carts := handler.CartService.GetByPhone(phone)
		res.Json(w, carts, http.StatusOK)
	}
}

func (handler *CartHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение данных
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, er.ErrNotAuthorized, http.StatusUnauthorized)
		}

		body, err := req.HandleBody[OrderRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Обновление в базе данных
		updatedCart, err := handler.CartService.Update(id, &Cart{
			Model:    gorm.Model{ID: uint(id)},
			Phone:    phone,
			Products: body.Products,
			Date:     datatypes.Date(time.Now()),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, updatedCart, http.StatusOK)
	}
}

func (handler *CartHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение данных
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, er.ErrNotAuthorized, http.StatusUnauthorized)
		}

		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Удаление данных
		deletedCart, err := handler.CartService.Delete(id, phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, deletedCart, http.StatusAccepted)
	}
}
