package main

import (
	"context"
	"time"

	"qasynda/shared/pkg/auth"
	"qasynda/shared/pkg/logger"
	pb "qasynda/shared/proto"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	store     *UserStore
	jwtSecret string
}

func NewServer(store *UserStore, jwtSecret string) *Server {
	return &Server{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	existing, err := s.store.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("failed to check existing user", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	if existing != nil {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	userID := uuid.New()
	user := &User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		FullName:     req.FullName,
		Phone:        req.Phone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.store.Create(ctx, user); err != nil {
		logger.Error("failed to create user", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	token, err := auth.GenerateToken(
		user.ID.String(),
		user.Email,
		user.Role,
		s.jwtSecret,
		24*time.Hour,
	)
	if err != nil {
		logger.Error("failed to generate token", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.AuthResponse{
		Token: token,
		User: &pb.UserResponse{
			Id:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			Phone:    user.Phone,
		},
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, err := s.store.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("failed to get user", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := auth.GenerateToken(
		user.ID.String(),
		user.Email,
		user.Role,
		s.jwtSecret,
		24*time.Hour,
	)
	if err != nil {
		logger.Error("failed to generate token", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.AuthResponse{
		Token: token,
		User: &pb.UserResponse{
			Id:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			Phone:    user.Phone,
		},
	}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.UserResponse, error) {
	claims, err := auth.ValidateToken(req.Token, s.jwtSecret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	uid, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token claims")
	}

	user, err := s.store.GetByID(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.UserResponse{
		Id:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		Phone:    user.Phone,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	user, err := s.store.GetByID(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.UserResponse{
		Id:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		Phone:    user.Phone,
	}, nil
}

func (s *Server) ListProviders(ctx context.Context, req *pb.ListProvidersRequest) (*pb.ListProvidersResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	providers, err := s.store.ListProviders(limit, offset)
	if err != nil {
		logger.Error("failed to list providers", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	var pbProviders []*pb.ProviderResponse
	for _, p := range providers {
		pbProviders = append(pbProviders, &pb.ProviderResponse{
			User: &pb.UserResponse{
				Id:       p.ID.String(),
				Email:    p.Email,
				FullName: p.FullName,
				Role:     p.Role,
				Phone:    p.Phone,
			},
			Location:        p.Location,
			HourlyRate:      p.HourlyRate,
			ExperienceYears: p.ExperienceYears,
			Bio:             p.Bio,
			IsAvailable:     p.IsAvailable,
			Rating:          p.Rating,
			ProviderId:      p.ServiceProviderID.String(),
		})
	}

	return &pb.ListProvidersResponse{Providers: pbProviders}, nil
}
