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
	ErrReviewNotFound     = errors.New("review not found")
	ErrBookingNotCompleted = errors.New("booking must be completed before review")
	ErrReviewExists       = errors.New("review already exists for this booking")
)

type ReviewService interface {
	Create(ctx context.Context, userID uuid.UUID, req *models.CreateReviewRequest) (*models.Review, error)
	GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.Review, error)
}

type reviewService struct {
	reviewRepo   repository.ReviewRepository
	bookingRepo  repository.BookingRepository
	providerRepo repository.ProviderRepository
	userRepo     repository.UserRepository
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	bookingRepo repository.BookingRepository,
	providerRepo repository.ProviderRepository,
	userRepo repository.UserRepository,
) ReviewService {
	return &reviewService{
		reviewRepo:   reviewRepo,
		bookingRepo:  bookingRepo,
		providerRepo: providerRepo,
		userRepo:     userRepo,
	}
}

func (s *reviewService) Create(ctx context.Context, userID uuid.UUID, req *models.CreateReviewRequest) (*models.Review, error) {
	// Verify booking exists and belongs to user
	booking, err := s.bookingRepo.GetByID(ctx, req.BookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	if booking.ClientID != userID {
		return nil, ErrUnauthorized
	}

	// Check if booking is completed
	if booking.Status != models.StatusCompleted {
		return nil, ErrBookingNotCompleted
	}

	// Check if review already exists
	existing, _ := s.reviewRepo.GetByBookingID(ctx, req.BookingID)
	if existing != nil {
		return nil, ErrReviewExists
	}

	review := &models.Review{
		ID:         uuid.New(),
		BookingID:  req.BookingID,
		ClientID:   userID,
		ProviderID: booking.ProviderID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		CreatedAt:  time.Now(),
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	// Update provider rating (this could be done by background worker)
	avgRating, totalReviews, _ := s.reviewRepo.GetProviderRating(ctx, booking.ProviderID)
	provider, _ := s.providerRepo.GetByID(ctx, booking.ProviderID)
	if provider != nil {
		provider.Rating = avgRating
		provider.TotalReviews = totalReviews
		s.providerRepo.Update(ctx, provider)
	}

	return review, nil
}

func (s *reviewService) GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.Review, error) {
	reviews, err := s.reviewRepo.GetByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	// Load client info
	for _, review := range reviews {
		client, _ := s.userRepo.GetByID(ctx, review.ClientID)
		review.Client = client
	}

	return reviews, nil
}

