package middlewares

import (
	"net/http"
)

const SecretKey = "secret"

func GenerateJWT() {}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)

		next.ServeHTTP(w, r)
	})
}
