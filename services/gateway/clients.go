package main

import (
	"log"

	"qasynda/shared/pkg/config"
	pb "qasynda/shared/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	User        pb.UserServiceClient
	Marketplace pb.MarketplaceServiceClient
	Chat        pb.ChatServiceClient
}

func InitClients(cfg *config.Config) *Clients {
	userConn, err := grpc.Dial(cfg.Services.UserUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to user service: %v", err)
	}

	marketplaceConn, err := grpc.Dial(cfg.Services.MarketplaceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to marketplace service: %v", err)
	}

	chatConn, err := grpc.Dial(cfg.Services.ChatUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to chat service: %v", err)
	}

	return &Clients{
		User:        pb.NewUserServiceClient(userConn),
		Marketplace: pb.NewMarketplaceServiceClient(marketplaceConn),
		Chat:        pb.NewChatServiceClient(chatConn),
	}
}
