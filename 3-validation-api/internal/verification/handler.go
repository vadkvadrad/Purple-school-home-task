package verification

import (
	"net/http"
)

type VerificationHandler struct {
	VerificationService *VerificationService
}

type VerificationHandlerDeps struct {
	VerificationService *VerificationService
}

func NewVerificationHandler(router *http.ServeMux, deps VerificationHandlerDeps) {
	handler := &VerificationHandler{
		deps.VerificationService,
	}

	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("POST /verify{hash}", handler.Verify())
}


func (handler *VerificationHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}


func (handler *VerificationHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

