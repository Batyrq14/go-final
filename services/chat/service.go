package main

import (
	"context"

	"qasynda/shared/pkg/logger"
	pb "qasynda/shared/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedChatServiceServer
	store *Store
}

func NewServer(store *Store) *Server {
	return &Server{store: store}
}

func (s *Server) GetHistory(ctx context.Context, req *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	u1, err := uuid.Parse(req.UserId_1)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id 1")
	}
	u2, err := uuid.Parse(req.UserId_2)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id 2")
	}

	messages, err := s.store.GetHistory(ctx, u1, u2, int(req.Limit), int(req.Offset))
	if err != nil {
		logger.Error("failed to get history", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	var pbMessages []*pb.Message
	for _, m := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			Id:         m.ID.String(),
			SenderId:   m.SenderID.String(),
			ReceiverId: m.ReceiverID.String(),
			Content:    m.Content,
			Timestamp:  m.CreatedAt.String(),
		})
	}

	return &pb.GetHistoryResponse{
		Messages: pbMessages,
	}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// This RPC is optional if we use WS, but useful for bots or system messages.
	// For now, implement basic save.
	return nil, status.Error(codes.Unimplemented, "use websocket for sending messages")
}
