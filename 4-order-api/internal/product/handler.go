package product

import (
	"net/http"
	"order-api/configs"
	"order-api/pkg/res"
)

type ProductHandlerDeps struct {
	Config *configs.Config
}

type ProductHandler struct {

}

func NewProductHandler(router *http.ServeMux, deps *ProductHandlerDeps) {
	handler := &ProductHandler{

	}

	router.HandleFunc("POST /product", handler.Create())
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res.Json(w, "done good", http.StatusOK)
	}
}