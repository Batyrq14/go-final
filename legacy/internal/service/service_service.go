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
	ErrServiceNotFound = errors.New("service not found")
)

type ServiceService interface {
	Create(ctx context.Context, req *models.CreateServiceRequest) (*models.Service, error)
	List(ctx context.Context) ([]*models.Service, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Service, error)
}

type serviceService struct {
	serviceRepo repository.ServiceRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository) ServiceService {
	return &serviceService{
		serviceRepo: serviceRepo,
	}
}

func (s *serviceService) Create(ctx context.Context, req *models.CreateServiceRequest) (*models.Service, error) {
	service := &models.Service{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.serviceRepo.Create(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *serviceService) List(ctx context.Context) ([]*models.Service, error) {
	return s.serviceRepo.List(ctx)
}

func (s *serviceService) GetByID(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	service, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrServiceNotFound
	}
	return service, nil
}

