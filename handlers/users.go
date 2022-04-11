package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

	if u.BlockedUntil.After(time.Now()) {
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(fmt.Sprintf("account locked, try again in %s.", time.Until(u.BlockedUntil).Round(time.Second)))
		return
	}

	if err := u.ComparePassword(body["password"]); err != nil {
		u.FailedLogins++

		if u.FailedLogins >= 5 {
			u.BlockedUntil = time.Now().Add(time.Second * 3)
		}

		h.store.UpdateOneUser(u)

		h.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u.FailedLogins = 0

	h.store.UpdateOneUser(u)

	th := NewTokensHandler(h.store, h.logger)

	token, err := th.GenerateJWT(u.Email)
	if err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expires := time.Now().Add(time.Hour * 2160)

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

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Unix(0, 0),
	})
}

func (h *UsersHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(r.Context().Value(CtxUserKey{}))
}

func (h *UsersHandler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	var body map[string]string

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u := r.Context().Value(CtxUserKey{}).(data.User)

	if body["first_name"] != "" {
		u.FirstName = body["first_name"]
	}

	if body["last_name"] != "" {
		u.LastName = body["last_name"]
	}

	if err := h.store.UpdateOneUser(u); err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *UsersHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var body map[string]string

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u := r.Context().Value(CtxUserKey{}).(data.User)

	u.SetPassword(body["password"])

	if err := h.store.UpdateOneUser(u); err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

		// add the user to the context
		ctx := context.WithValue(r.Context(), KeyBody{}, body)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
