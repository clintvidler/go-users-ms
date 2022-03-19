package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"users-ms/data"
	"users-ms/middlewares"
)

type UsersHandler struct {
	store  *data.Store
	logger *middlewares.Logger
}

func NewUsersHandler(ds *data.Store, l *middlewares.Logger) *UsersHandler {
	return &UsersHandler{store: ds, logger: l}
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	body := r.Context().Value(KeyBody{}).(map[string]string)

	u, err := h.store.ReadOneUser(body["email"])
	if err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if u.Email == "" {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := u.ComparePassword(body["password"]); err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	th := NewTokensHandler(h.store, h.logger)

	token, err := th.GenerateJWT(u.Email)
	if err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// data activity to login
	expires := time.Now().Add(time.Minute * 5)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expires,
	})

	h.store.Login(u, token, expires)

	json.NewEncoder(w).Encode(token)
}

func (h *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	body := r.Context().Value(KeyBody{}).(map[string]string)

	user := data.User{
		FirstName: body["first_name"],
		LastName:  body["last_name"],
		Email:     body["email"],
	}

	if err := user.SetPassword(body["password"]); err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.store.CreateOneUser(&user); err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user.Email)
}

func (h *UsersHandler) Logout(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(CtxUserKey{})

	h.store.DeleteOneToken(u.(data.User).Email)
}

type KeyBody struct{}

func (h *UsersHandler) ValidateInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			h.logger.Debug(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := data.User{
			FirstName: body["first_name"],
			LastName:  body["last_name"],
			Email:     body["email"],
		}

		if err := user.SetPassword(body["password"]); err != nil {
			h.logger.Debug(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := user.Validate(); err != nil {
			h.logger.Debug(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyBody{}, body)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
