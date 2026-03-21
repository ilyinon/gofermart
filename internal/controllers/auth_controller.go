package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"gophermart/internal/services"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("register decode error", "err", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}

	err := c.service.Register(r.Context(), req.Login, req.Password)
	if err != nil {

		if errors.Is(err, services.ErrUserExists) {
			http.Error(w, "login already exists", http.StatusConflict)
			return
		}

		slog.Error("register error", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, err := c.service.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		slog.Error("auto login error", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("login decode error", "err", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}

	token, err := c.service.Login(r.Context(), req.Login, req.Password)
	if err != nil {

		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	w.WriteHeader(http.StatusOK)
}
