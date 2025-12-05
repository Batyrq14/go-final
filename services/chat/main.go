package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"qasynda/shared/pkg/config"
	"qasynda/shared/pkg/db"
	"qasynda/shared/pkg/logger"
	pb "qasynda/shared/proto"

	"google.golang.org/grpc"
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
		// Non-fatal for demo? Fatal is better.
		os.Exit(1)
	}
	defer rmq.Close()

	// Start consumer
	go StartConsumer(cfg.RabbitMQUrl, store)

	// Start WS Hub
	hub := NewHub(store, rmq)
	go hub.Run()

	// WS Handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	// GRPC Server
	port := config.GetChatPort()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, NewServer(store))

	logger.Info("Chat Service (gRPC) starting on " + port)

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Error("failed to serve grpc", err)
			os.Exit(1)
		}
	}()

	// HTTP Server for WS
	// Use 8081 or verify if we can share port. We cannot share port easily.
	wsPort := ":8081"
	logger.Info("Chat Service (WS) starting on " + wsPort)
	go func() {
		if err := http.ListenAndServe(wsPort, nil); err != nil {
			logger.Error("failed to serve http", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Chat Service...")
	s.GracefulStop()
}
