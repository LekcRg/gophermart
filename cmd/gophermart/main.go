package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/LekcRg/gophermart/docs"
	"github.com/LekcRg/gophermart/internal/accrual"
	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/handlers"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/repository/postgres"
	"github.com/LekcRg/gophermart/internal/router"
	"github.com/LekcRg/gophermart/internal/service"
	"github.com/LekcRg/gophermart/internal/validator"
	"go.uber.org/zap"
)

// @title Gophermart API
// @version 1.0
// @description Gophermart cumulative loyalty system

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func handleGracefulShutdown(
	cancel context.CancelFunc, server *http.Server,
	db repository.RepositoryProvider, acr *accrual.Accrual,
) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Log.Info("Shutting down")
	cancel()

	logger.Log.Info("stopping db")
	db.Close()

	logger.Log.Info("stopping server")
	ctx, toCancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Log.Error("error while shutdown server", zap.Error(err))
	}

	err = logger.Log.Sync()
	if err != nil {
		logger.Log.Error("Log.Sync error", zap.Error(err))
	}

	acr.Close()
	toCancel()
}

func main() {
	cfg := config.Get()
	logger.Initialize(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	dbProvider := postgres.New(ctx, cfg)
	valid := validator.New()
	repos := dbProvider.GetRepositories()
	req := accrual.New(cfg.AccrualAddress)
	services := service.New(ctx, repos, valid, cfg, req)
	handlers := handlers.New(cfg, services, valid)
	routes := router.New(handlers, cfg.JWTSecret)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: routes,
	}
	go handleGracefulShutdown(cancel, server, dbProvider, req)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error(err.Error())
	}
	logger.Log.Info("Buy, 👋!")
}
