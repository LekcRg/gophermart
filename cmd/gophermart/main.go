package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/router"
	"github.com/LekcRg/gophermart/internal/storage"
	"github.com/LekcRg/gophermart/internal/storage/postgres"
	"go.uber.org/zap"
)

func exit(cancel context.CancelFunc, server *http.Server, db storage.Storage) {
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
	routes := router.New()
	logger.Log.Info("Hello, world!")
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: routes,
	}

	ctx, cancel := context.WithCancel(context.Background())
	db := postgres.New(ctx, cfg)

	go exit(cancel, server, db)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error(err.Error())
	}
	logger.Log.Info("Buy, 👋!")
}
