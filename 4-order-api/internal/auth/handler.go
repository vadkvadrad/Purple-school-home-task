package auth

import (
	"net/http"
	"order-api/configs"
	"order-api/pkg/jwt"
	"order-api/pkg/req"
	"order-api/pkg/res"
)

type AuthHandler struct {
	AuthService *AuthService
	Config *configs.Config
}

type AuthHandlerDeps struct {
	AuthService *AuthService
	Config *configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		AuthService: deps.AuthService,
		Config: deps.Config,
	}

	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/verify", handler.Verify())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sessionId, err := handler.AuthService.Login(body.Phone, body.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := LoginResponse{
			SessionId: sessionId,
		}
		res.Json(w, data, http.StatusCreated)
	}
}

func (handler *AuthHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[VerifyRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := handler.AuthService.GetBySessionId(body.SessionId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Code != body.Code {
			res.Json(w, "wrong authorization code", http.StatusUnauthorized)
			return 
		}

		jwtCreator := jwt.NewJwt(handler.Config.Auth.Secret)
		token, err := jwtCreator.Create(jwt.JWTData{
			Phone: user.Phone,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		data := &VerifyResponse{
			Token: token,
		} 
		res.Json(w, data, http.StatusOK)
	}
}