package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gophermart/internal/config"
	"gophermart/internal/controllers"
	"gophermart/internal/infrastructure"
	"gophermart/internal/repositories"
	"gophermart/internal/services"
	"gophermart/internal/utils"
	"gophermart/internal/worker"
)

func main() {

	cfg := config.Load()

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	).With(
		"service", "gophermart",
		"env", "dev",
	)

	slog.SetDefault(logger)

	db, err := infrastructure.NewPostgres(cfg.DatabaseURI)
	if err != nil {
		slog.Error("database connection error", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()

	if err := infrastructure.EnsureSchema(ctx, db); err != nil {
		slog.Error("schema init failed", "err", err)
		os.Exit(1)
	}

	userRepo := repositories.NewUserRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	withdrawRepo := repositories.NewWithdrawalRepository(db)

	jwtManager := utils.NewJWTManager(cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, jwtManager)

	orderService := services.NewOrderService(orderRepo)
	balanceService := services.NewBalanceService(orderRepo, withdrawRepo)

	authController := controllers.NewAuthController(authService)
	orderController := controllers.NewOrderController(orderService)
	balanceController := controllers.NewBalanceController(balanceService)

	router := controllers.NewRouter(
		authController,
		orderController,
		balanceController,
		jwtManager,
	)
	ctxWorker, cancelWorker := context.WithCancel(context.Background())

	accrualClient := infrastructure.NewAccrualClient(cfg.AccrualSystemAddress)

	accrualWorker := worker.NewAccrualWorker(
		orderRepo,
		accrualClient,
		ctxWorker,
	)

	go accrualWorker.Start(ctxWorker)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		slog.Info("server started", "addr", cfg.RunAddress)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	<-stop

	slog.Info("shutdown initiated")

	cancelWorker()

	ctxShutdown, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		slog.Error("server shutdown error", "err", err)
	}

	accrualWorker.Stop()

	db.Close()

	slog.Info("server stopped")
}
