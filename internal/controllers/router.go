package controllers

import (
	"gophermart/internal/middleware"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	auth *AuthController,
	order *OrderController,
	balance *BalanceController,
) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Post("/api/user/register", auth.Register)
	r.Post("/api/user/login", auth.Login)

	r.Group(func(r chi.Router) {

		r.Use(middleware.Auth)

		r.Post("/api/user/orders", order.Upload)
		r.Get("/api/user/orders", order.List)

		r.Get("/api/user/balance", balance.GetBalance)

		r.Post("/api/user/balance/withdraw", balance.Withdraw)

		r.Get("/api/user/withdrawals", balance.Withdrawals)
	})

	return r
}
