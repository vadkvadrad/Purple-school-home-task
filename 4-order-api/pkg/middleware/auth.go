package middleware

import (
	"context"
	"net/http"
	"order-api/configs"
	"order-api/pkg/jwt"
	"strings"
)

type key string

const (
	ContextPhoneKey key = "ContextPhoneKey"
	ErrUnauthorized string = "unauthorized"
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authedHeader, "Bearer ") {
			writeUnauthed(w)
			return
		}
		token := strings.TrimPrefix(authedHeader, "Bearer ")
		isValid, data := jwt.NewJwt(config.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthed(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextPhoneKey, data.Phone)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	http.Error(w, ErrUnauthorized, http.StatusBadRequest)
}