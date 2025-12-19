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

	rmq, err := NewRabbitMQProducer(cfg.RabbitMQUrl)
	if err != nil {
		logger.Error("failed to connect to rabbitmq", err)
		os.Exit(1)
	}
	defer rmq.Close()

	ctx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	go StartConsumer(ctx, cfg.RabbitMQUrl, store)

	hub := NewHub(store, rmq)
	go hub.Run(ctx)

	server := NewServer(store)

	r := gin.Default()

	r.GET("/history", server.GetHistory)
	r.GET("/ws", func(c *gin.Context) {
		ServeWs(hub, c.Writer, c.Request)
	})

	port := config.GetChatPort()
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		logger.Info("Chat Service starting HTTP/WS on " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to serve", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down Chat Service...")

	cancelWorkers()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Chat Service forced to shutdown", err)
	}

	logger.Info("Chat Service exiting")
}
