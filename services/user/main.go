package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qasynda/shared/pkg/config"
	"qasynda/shared/pkg/db"
	"qasynda/shared/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()
	cfg := config.Load()

	database, err := db.Connect(cfg.DBUrl)
	if err != nil {
		logger.Error("failed to connect to db", err)
		os.Exit(1)
	}
	defer database.Close()

	store := NewUserStore(database)
	server := NewServer(store, cfg.JWTSecret)

	r := gin.Default()

	r.POST("/register", server.Register)
	r.POST("/login", server.Login)
	r.POST("/validate", server.ValidateToken)
	r.GET("/users/:id", server.GetUser)
	r.GET("/providers", server.ListProviders)
	r.PUT("/providers/:id/status", server.UpdateProviderStatus)
	r.GET("/providers/:id/status", server.GetProviderStatus)

	port := config.GetUserPort()
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		logger.Info("User Service starting HTTP on " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to serve", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down User Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("User Service forced to shutdown", err)
	}

	logger.Info("User Service exiting")
}
