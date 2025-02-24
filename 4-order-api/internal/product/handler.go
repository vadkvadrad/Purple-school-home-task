package product

import (
	"net/http"
	"order-api/configs"
	"order-api/pkg/er"
	"order-api/pkg/middleware"
	"order-api/pkg/req"
	"order-api/pkg/res"
	"strconv"
)

const (
	ErrWrongProductCredentials = "wrong user credentials"
)

type ProductHandler struct {
	Config *configs.Config
	ProductRepository *ProductRepository
}

type ProductHandlerDeps struct {
	Config      *configs.Config
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		Config:      deps.Config,
		ProductRepository: deps.ProductRepository,
	}

	router.Handle("POST /product", middleware.IsAuthed(handler.Create(), handler.Config))
	router.Handle("UPDATE /product/{id}", middleware.IsAuthed(handler.Update(), handler.Config))
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

		product, err := handler.ProductRepository.Create(&Product{
			Name: body.Name,
			Description: body.Description,
			Images: body.Images,
			Price: body.Price,
			Currency: CurrencyRUB,
			Owner: phone,
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

		prod, err := handler.ProductRepository.FindById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// service logic 
		if prod.Owner != phone {
			http.Error(w, er.ErrWrongUserCredentials, http.StatusBadRequest)
			return
		}

		prod, err = handler.ProductRepository.Update(&Product{
			Name: body.Name,
			Description: body.Description,
			Images: body.Images,
			Price: body.Price,
			Currency: CurrencyRUB,
			Owner: phone,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, prod, http.StatusOK)
	}
}
