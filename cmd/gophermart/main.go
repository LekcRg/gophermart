package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func exit(
	cancel context.CancelFunc, server *http.Server,
	db repository.RepositoryProvider,
) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Log.Info("Shutting down")
	cancel()

	logger.Log.Info("stopping db")
	db.Close()

	logger.Log.Info("stopping server")
	ctx, toCancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Log.Error("error while shutdown server", zap.Error(err))
	}
	toCancel()
}

func main() {
	cfg := config.Get()
	logger.Initialize(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	dbProvider := postgres.New(ctx, cfg)
	valid := validator.New()
	repos := dbProvider.GetRepositories()
	services := service.New(repos, valid, cfg)
	handlers := handlers.New(cfg, services, valid)
	routes := router.New(handlers, cfg.JWTSecret)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: routes,
	}
	go exit(cancel, server, dbProvider)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error(err.Error())
	}
	logger.Log.Info("Buy, 👋!")
}
