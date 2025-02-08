package verify

import (
	"net/http"
	"verify-api/pkg/req"
)

const (
	ErrUserAlreadyVerified = "user already verified"
)

type VerifyHandler struct {
	VerifyService *VerifyService
}

type VerifyHandlerDeps struct {
	VerifyService *VerifyService
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		deps.VerifyService,
	}

	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())

	router.HandleFunc("POST /test", handler.TmpSend())
}

func (handler *VerifyHandler) TmpSend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler.VerifyService.TemSend()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[SendRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// FIXME START
		currentUser, _ := handler.VerifyService.GetByEmail(body.Email)
		if currentUser != nil && currentUser.IsVerified {
			http.Error(w, ErrUserAlreadyVerified, http.StatusBadRequest)
			return
		} else if currentUser != nil && !currentUser.IsVerified {
			err = handler.VerifyService.Send(currentUser.Hash)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		// END


		hash := handler.VerifyService.GenerateHash()

		handler.VerifyService.Create(
			body.Email,
			hash,
		)

		err = handler.VerifyService.Send(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}


func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		user, err := handler.VerifyService.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if user.IsVerified {
			http.Error(w, ErrUserAlreadyVerified, http.StatusBadRequest)
			return
		}

		_, err = handler.VerifyService.Verify(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

