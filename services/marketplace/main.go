package main

import (
	"os"

	"qasynda/shared/pkg/config"
	"qasynda/shared/pkg/db"
	"qasynda/shared/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()
	cfg := config.Load()

	// Connect to DB
	database, err := db.Connect(cfg.DBUrl)
	if err != nil {
		logger.Error("failed to connect to db", err)
		os.Exit(1)
	}
	defer database.Close()

	// Init Store
	store := NewStore(database)
	server := NewServer(store)

	// Init Gin
	r := gin.Default()

	// Define Routes
	r.GET("/services", server.GetServices)
	r.POST("/services", server.CreateService)
	r.GET("/bookings", server.ListBookings)
	r.POST("/bookings", server.CreateBooking)
	r.PUT("/bookings/:id/status", server.UpdateBookingStatus)

	port := config.GetMarketplacePort() // e.g. :50052
	logger.Info("Marketplace Service starting HTTP on " + port)

	if err := r.Run(port); err != nil {
		logger.Error("failed to serve", err)
		os.Exit(1)
	}
}
