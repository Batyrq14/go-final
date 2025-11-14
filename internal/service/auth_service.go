package service

import (
	"context"
	"errors"
	"qasynda/internal/config"
	"qasynda/internal/models"
	"qasynda/internal/repository"
	"qasynda/pkg/auth"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	providerRepo repository.ProviderRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, providerRepo repository.ProviderRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		providerRepo: providerRepo,
		config:   cfg,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if email already exists
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		FullName:     req.FullName,
		Phone:        req.Phone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// If provider, create provider profile
	if req.Role == models.RoleProvider {
		provider := &models.ServiceProvider{
			ID:              uuid.New(),
			UserID:          user.ID,
			HourlyRate:      0,
			ExperienceYears: 0,
			IsAvailable:     true,
			Rating:          0,
			TotalReviews:    0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if err := s.providerRepo.Create(ctx, provider); err != nil {
			return nil, err
		}
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *authService) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *authService) generateToken(user *models.User) (string, error) {
	return auth.GenerateToken(
		user.ID.String(),
		user.Email,
		string(user.Role),
		s.config.JWT.Secret,
		s.config.JWT.Expiration,
	)
}

