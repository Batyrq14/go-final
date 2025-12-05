package service

import (
	"context"
	"errors"
	"qasynda/internal/models"
	"qasynda/internal/repository"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrUnauthorized    = errors.New("unauthorized")
)

type ProviderService interface {
	List(ctx context.Context, filter *models.ProviderFilter) (*models.ProviderListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.ServiceProvider, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdateProviderRequest) (*models.ServiceProvider, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type providerService struct {
	providerRepo repository.ProviderRepository
	userRepo     repository.UserRepository
	serviceRepo  repository.ServiceRepository
}

func NewProviderService(providerRepo repository.ProviderRepository, userRepo repository.UserRepository, serviceRepo repository.ServiceRepository) ProviderService {
	return &providerService{
		providerRepo: providerRepo,
		userRepo:     userRepo,
		serviceRepo:  serviceRepo,
	}
}

func (s *providerService) List(ctx context.Context, filter *models.ProviderFilter) (*models.ProviderListResponse, error) {
	providers, total, err := s.providerRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Load user and services for each provider
	for _, provider := range providers {
		user, _ := s.userRepo.GetByID(ctx, provider.UserID)
		provider.User = user

		services, _ := s.providerRepo.GetProviderServices(ctx, provider.ID)
		provider.Services = services
	}

	return &models.ProviderListResponse{
		Data: providers,
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
			Total: total,
		},
	}, nil
}

func (s *providerService) GetByID(ctx context.Context, id uuid.UUID) (*models.ServiceProvider, error) {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProviderNotFound
	}

	// Load user
	user, _ := s.userRepo.GetByID(ctx, provider.UserID)
	provider.User = user

	// Load services
	services, _ := s.providerRepo.GetProviderServices(ctx, provider.ID)
	provider.Services = services

	return provider, nil
}

func (s *providerService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdateProviderRequest) (*models.ServiceProvider, error) {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProviderNotFound
	}

	// Check authorization
	if provider.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields
	if req.Bio != nil {
		provider.Bio = *req.Bio
	}
	if req.HourlyRate != nil {
		provider.HourlyRate = *req.HourlyRate
	}
	if req.ExperienceYears != nil {
		provider.ExperienceYears = *req.ExperienceYears
	}
	if req.Location != nil {
		provider.Location = *req.Location
	}
	if req.IsAvailable != nil {
		provider.IsAvailable = *req.IsAvailable
	}
	provider.UpdatedAt = time.Now()

	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *providerService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return ErrProviderNotFound
	}

	// Check authorization
	if provider.UserID != userID {
		return ErrUnauthorized
	}

	return s.providerRepo.Delete(ctx, id)
}

