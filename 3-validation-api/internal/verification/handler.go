package verification

import (
	"net/http"
)

type EmailHandler struct {
	VerificationService *VerificationService
}

type EmailHandlerDeps struct {
	VerificationService *VerificationService
}

func NewEmailHandler(router *http.ServeMux, deps EmailHandlerDeps) {
	handler := &EmailHandler{
		deps.VerificationService,
	}

	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("POST /verify{hash}", handler.Verify())
}


func (handler *EmailHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}


func (handler *EmailHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

