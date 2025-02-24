package product

import (
	"net/http"
	"order-api/configs"
	"order-api/pkg/req"
	"order-api/pkg/res"
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

	router.HandleFunc("POST /product", handler.Create())
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		res.Json(w, product, http.StatusCreated)
	}
}
