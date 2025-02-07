package verification

import (
	"net/http"
	"verify-api/configs"
)

type EmailHandler struct {
	Config *configs.Config
}

type EmailHandlerDeps struct {
	// service here
}

func NewEmailHandler(router *http.ServeMux, deps EmailHandlerDeps) {
	handler := &EmailHandler{
		// service here
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

