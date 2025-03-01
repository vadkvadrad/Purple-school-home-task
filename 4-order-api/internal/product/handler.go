package product

import (
	"net/http"
	"order-api/configs"
	"order-api/pkg/er"
	"order-api/pkg/middleware"
	"order-api/pkg/req"
	"order-api/pkg/res"
	"strconv"

	"gorm.io/gorm"
)

const (
	ErrWrongProductCredentials = "wrong user credentials"
)

type ProductHandler struct {
	Config         *configs.Config
	ProductService *ProductService
}

type ProductHandlerDeps struct {
	Config         *configs.Config
	ProductService *ProductService
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		Config:         deps.Config,
		ProductService: deps.ProductService,
	}

	// Добавить продукт
	router.Handle("POST /product", middleware.IsAuthed(handler.Create(), handler.Config))

	// Обновить продукт
	router.Handle("PATCH /product/{id}", middleware.IsAuthed(handler.Update(), handler.Config))

	// Удалить продукт
	router.Handle("DELETE /product/{id}", middleware.IsAuthed(handler.Delete(), handler.Config))
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, er.ErrNotAuthorized, http.StatusUnauthorized)
			return
		}

		body, err := req.HandleBody[ProductRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		product, err := handler.ProductService.Create(&Product{
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
			Price:       body.Price,
			Currency:    CurrencyRUB,
			Owner:       phone,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, product, http.StatusCreated)
	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение данных
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, er.ErrNotAuthorized, http.StatusUnauthorized)
			return
		}

		body, err := req.HandleBody[ProductRequest](w, r)
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


		// Получение продукта
		prod, err := handler.ProductService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}


		// Обновление продукта
		prod, err = handler.ProductService.Update(prod.Owner, &Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
			Price:       body.Price,
			Currency:    CurrencyRUB,
			Owner:       phone,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, prod, http.StatusOK)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение данных
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


		// Получение продукта
		prod, err := handler.ProductService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		

		// Удаление продукта
		err = handler.ProductService.Delete(prod.Owner, phone, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, prod, http.StatusAccepted)
	}
}
