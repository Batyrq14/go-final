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
	store := NewStore(database)

	// Init GRPC Server
	port := config.GetMarketplacePort()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterMarketplaceServiceServer(s, NewServer(store))

	logger.Info("Marketplace Service starting on " + port)

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

	logger.Info("Shutting down Marketplace Service...")
	s.GracefulStop()
}
