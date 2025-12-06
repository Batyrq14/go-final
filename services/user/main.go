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
	store := NewUserStore(database)
	server := NewServer(store, cfg.JWTSecret)

	// Init Gin
	r := gin.Default()

	// Define Routes
	r.POST("/register", server.Register)
	r.POST("/login", server.Login)
	r.POST("/validate", server.ValidateToken)
	r.GET("/users/:id", server.GetUser)
	r.GET("/providers", server.ListProviders)

	port := config.GetUserPort() // e.g. :50051 (we might want to change this to a standard HTTP port like :8081 eventually, but keeping config port is fine for now)
	logger.Info("User Service starting HTTP on " + port)

	if err := r.Run(port); err != nil {
		logger.Error("failed to serve", err)
		os.Exit(1)
	}
}
