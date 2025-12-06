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

	store := NewStore(database)

	// RabbitMQ
	rmq, err := NewRabbitMQProducer(cfg.RabbitMQUrl)
	if err != nil {
		logger.Error("failed to connect to rabbitmq", err)
		os.Exit(1)
	}
	defer rmq.Close()

	// Start consumer
	go StartConsumer(cfg.RabbitMQUrl, store)

	// Start WS Hub
	hub := NewHub(store, rmq)
	go hub.Run()

	// Init Server
	server := NewServer(store)

	// Init Gin
	r := gin.Default()

	// Define Routes
	r.GET("/history", server.GetHistory)
	r.GET("/ws", func(c *gin.Context) {
		ServeWs(hub, c.Writer, c.Request)
	})

	// Use one port for both HTTP Routes and WS
	port := config.GetChatPort() // e.g. :50053 (will be HTTP now)
	logger.Info("Chat Service starting HTTP/WS on " + port)

	if err := r.Run(port); err != nil {
		logger.Error("failed to serve", err)
		os.Exit(1)
	}
}
