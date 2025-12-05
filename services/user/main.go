package main

import (
	"net"
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

	// Init Store
	store := NewUserStore(database)

	// Init GRPC Server
	lis, err := net.Listen("tcp", cfg.Services.UserUrl) // Assuming UserUrl is like "localhost:50051" or ":50051". Actually config default was localhost:50051. Listen needs :port

	// Fix: Config returns full URL, but listen needs just port if running in container or host.
	// We'll use getEnv(":50051") sort of logic elsewhere or just use :50051 if running locally.
	// But docker-compose will network them.
	// Let's rely on config.GetUserPort() I added earlier.
	port := config.GetUserPort()
	lis, err = net.Listen("tcp", port)
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, NewServer(store, cfg.JWTSecret))

	logger.Info("User Service starting on " + port)

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Error("failed to serve", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down User Service...")
	s.GracefulStop()
}
