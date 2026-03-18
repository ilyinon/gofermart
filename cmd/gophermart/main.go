package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gophermart/internal/config"
	"gophermart/internal/controllers"
	"gophermart/internal/infrastructure"
	"gophermart/internal/repositories"
	"gophermart/internal/services"
	"gophermart/internal/worker"
)

func main() {

	cfg := config.Load()

	db, err := infrastructure.NewPostgres(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	ctx := context.Background()

	if err := infrastructure.EnsureSchema(ctx, db); err != nil {
		log.Fatalf("schema init failed: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	withdrawRepo := repositories.NewWithdrawalRepository(db)

	authService := services.NewAuthService(userRepo)
	orderService := services.NewOrderService(orderRepo)
	balanceService := services.NewBalanceService(orderRepo, withdrawRepo)

	authController := controllers.NewAuthController(authService)
	orderController := controllers.NewOrderController(orderService)
	balanceController := controllers.NewBalanceController(balanceService)

	router := controllers.NewRouter(
		authController,
		orderController,
		balanceController,
	)

	accrualClient := infrastructure.NewAccrualClient(cfg.AccrualSystemAddress)

	accrualWorker := worker.NewAccrualWorker(orderRepo, accrualClient)

	ctxWorker, cancelWorker := context.WithCancel(context.Background())

	go accrualWorker.Start(ctxWorker)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		log.Printf("server started on %s", cfg.RunAddress)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Println("shutdown initiated")

	cancelWorker()

	ctxShutdown, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	accrualWorker.Stop()

	db.Close()

	log.Println("server stopped")
}
