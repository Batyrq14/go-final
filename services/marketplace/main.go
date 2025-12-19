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

	store := NewStore(database)
	server := NewServer(store)

	r := gin.Default()

	r.GET("/services", server.GetServices)
	r.POST("/services", server.CreateService)
	r.GET("/bookings", server.ListBookings)
	r.POST("/bookings", server.CreateBooking)
	r.PUT("/bookings/:id/status", server.UpdateBookingStatus)

	port := config.GetMarketplacePort()
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		logger.Info("Marketplace Service starting HTTP on " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to serve", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down Marketplace Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Marketplace Service forced to shutdown", err)
	}

	logger.Info("Marketplace Service exiting")
}
