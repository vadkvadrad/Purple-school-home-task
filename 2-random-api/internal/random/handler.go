package random

import (
	"math/rand"
	"net/http"
	"rand-api/pkg/res"
//	"strconv"
)

type RandomHandler struct{}

func NewRandomHandler(router *http.ServeMux) {
	handler := &RandomHandler{}
	router.HandleFunc("GET /random", handler.Random())
}

func (handler *RandomHandler) Random() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte(strconv.Itoa(rand.Intn(5)+1)))
		res.Json(w, rand.Intn(6)+1, http.StatusOK)
	}
}