package handlers

import (
	"context"
	"net/http"
	"time"
	"users-ms/data"
	"users-ms/middlewares"

	"github.com/golang-jwt/jwt"
)

type TokensHandler struct {
	store  *data.Store
	logger *middlewares.Logger
	secret string
}

func NewTokensHandler(s *data.Store, l *middlewares.Logger) *TokensHandler {
	return &TokensHandler{store: s, logger: l, secret: "secret"}
}

func (h *TokensHandler) GenerateJWT(email string) (string, error) {
	payload := jwt.StandardClaims{}
	payload.Subject = email
	payload.ExpiresAt = time.Now().Add(time.Minute * 5).Unix()

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(h.secret))
}

type CtxUserKey struct{} // used as key in context k, v

func (h *TokensHandler) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ct, err := jwt.ParseWithClaims(c.Value, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(h.secret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !ct.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		payload := ct.Claims.(*jwt.StandardClaims)

		email := payload.Subject

		t, err := h.store.ReadOneToken(email, ct.Raw)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if t.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		u, err := h.store.ReadOneUser(email)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserKey{}, u)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
