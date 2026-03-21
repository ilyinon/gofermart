package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"gophermart/internal/middleware"
	"gophermart/internal/services"
	"gophermart/internal/utils"
)

type OrderController struct {
	service *services.OrderService
}

func NewOrderController(service *services.OrderService) *OrderController {
	return &OrderController{
		service: service,
	}
}

func (c *OrderController) Upload(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	number := strings.TrimSpace(string(body))

	if !utils.ValidLuhn(number) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = c.service.Upload(r.Context(), userID, number)

	if err != nil {

		slog.Error("order upload error:", "err", err)

		if errors.Is(err, services.ErrOrderExists) {
			w.WriteHeader(http.StatusOK)
			return
		}

		if errors.Is(err, services.ErrOrderUsed) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (c *OrderController) List(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := c.service.List(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(orders)
}
