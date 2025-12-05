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
	ErrBookingNotFound = errors.New("booking not found")
	ErrInvalidStatus   = errors.New("invalid status transition")
)

type BookingService interface {
	Create(ctx context.Context, userID uuid.UUID, req *models.CreateBookingRequest) (*models.Booking, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole models.UserRole) (*models.Booking, error)
	List(ctx context.Context, userID uuid.UUID, userRole models.UserRole) ([]*models.Booking, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole models.UserRole, req *models.UpdateBookingStatusRequest) (*models.Booking, error)
}

type bookingService struct {
	bookingRepo  repository.BookingRepository
	providerRepo repository.ProviderRepository
	serviceRepo  repository.ServiceRepository
	userRepo     repository.UserRepository
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	providerRepo repository.ProviderRepository,
	serviceRepo repository.ServiceRepository,
	userRepo repository.UserRepository,
) BookingService {
	return &bookingService{
		bookingRepo:  bookingRepo,
		providerRepo: providerRepo,
		serviceRepo:  serviceRepo,
		userRepo:     userRepo,
	}
}

func (s *bookingService) Create(ctx context.Context, userID uuid.UUID, req *models.CreateBookingRequest) (*models.Booking, error) {
	// Verify provider exists
	provider, err := s.providerRepo.GetByID(ctx, req.ProviderID)
	if err != nil {
		return nil, errors.New("provider not found")
	}

	// Verify service exists
	svc, err := s.serviceRepo.GetByID(ctx, req.ServiceID)
	if err != nil {
		return nil, errors.New("service not found")
	}
	_ = svc // Use service to verify it exists

	// Parse scheduled date
	scheduledDate, err := time.Parse(time.RFC3339, req.ScheduledDate)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// Calculate total price
	totalPrice := provider.HourlyRate * req.DurationHours

	booking := &models.Booking{
		ID:            uuid.New(),
		ClientID:      userID,
		ProviderID:    req.ProviderID,
		ServiceID:     req.ServiceID,
		ScheduledDate: scheduledDate,
		DurationHours: req.DurationHours,
		Status:        models.StatusPending,
		TotalPrice:    totalPrice,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole models.UserRole) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBookingNotFound
	}

	// Check authorization
	if userRole != models.RoleAdmin {
		if booking.ClientID != userID && booking.ProviderID != userID {
			return nil, ErrUnauthorized
		}
	}

	// Load related data
	client, _ := s.userRepo.GetByID(ctx, booking.ClientID)
	booking.Client = client

	provider, _ := s.providerRepo.GetByID(ctx, booking.ProviderID)
	booking.Provider = provider

	service, _ := s.serviceRepo.GetByID(ctx, booking.ServiceID)
	booking.Service = service

	return booking, nil
}

func (s *bookingService) List(ctx context.Context, userID uuid.UUID, userRole models.UserRole) ([]*models.Booking, error) {
	var bookings []*models.Booking
	var err error

	if userRole == models.RoleProvider {
		provider, _ := s.providerRepo.GetByUserID(ctx, userID)
		if provider != nil {
			bookings, err = s.bookingRepo.GetByProviderID(ctx, provider.ID)
		}
	} else {
		bookings, err = s.bookingRepo.GetByClientID(ctx, userID)
	}

	if err != nil {
		return nil, err
	}

	// Load related data
	for _, booking := range bookings {
		client, _ := s.userRepo.GetByID(ctx, booking.ClientID)
		booking.Client = client

		provider, _ := s.providerRepo.GetByID(ctx, booking.ProviderID)
		booking.Provider = provider

		service, _ := s.serviceRepo.GetByID(ctx, booking.ServiceID)
		booking.Service = service
	}

	return bookings, nil
}

func (s *bookingService) UpdateStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole models.UserRole, req *models.UpdateBookingStatusRequest) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBookingNotFound
	}

	// Check authorization and validate status transition
	if userRole == models.RoleProvider {
		provider, _ := s.providerRepo.GetByUserID(ctx, userID)
		if provider == nil || provider.ID != booking.ProviderID {
			return nil, ErrUnauthorized
		}
		// Providers can only accept or reject
		if req.Status != models.StatusAccepted && req.Status != models.StatusRejected {
			return nil, ErrInvalidStatus
		}
	} else if userRole == models.RoleClient {
		if booking.ClientID != userID {
			return nil, ErrUnauthorized
		}
		// Clients can only cancel
		if req.Status != models.StatusCancelled {
			return nil, ErrInvalidStatus
		}
	}

	// Validate status transition
	if booking.Status == models.StatusCompleted || booking.Status == models.StatusCancelled {
		return nil, errors.New("cannot change status of completed or cancelled booking")
	}

	if err := s.bookingRepo.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, err
	}

	booking.Status = req.Status
	booking.UpdatedAt = time.Now()

	return booking, nil
}

