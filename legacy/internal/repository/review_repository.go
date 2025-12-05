package repository

import (
	"context"
	"qasynda/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error)
	GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.Review, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetProviderRating(ctx context.Context, providerID uuid.UUID) (float64, int, error)
}

type reviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	query := `
		INSERT INTO reviews (id, booking_id, client_id, provider_id, rating, comment, created_at)
		VALUES (:id, :booking_id, :client_id, :provider_id, :rating, :comment, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, review)
	return err
}

func (r *reviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	var review models.Review
	query := `
		SELECT id, booking_id, client_id, provider_id, rating, comment, created_at
		FROM reviews WHERE id = $1
	`
	err := r.db.GetContext(ctx, &review, query, id)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.Review, error) {
	var reviews []*models.Review
	query := `
		SELECT id, booking_id, client_id, provider_id, rating, comment, created_at
		FROM reviews WHERE provider_id = $1 ORDER BY created_at DESC
	`
	err := r.db.SelectContext(ctx, &reviews, query, providerID)
	return reviews, err
}

func (r *reviewRepository) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*models.Review, error) {
	var review models.Review
	query := `
		SELECT id, booking_id, client_id, provider_id, rating, comment, created_at
		FROM reviews WHERE booking_id = $1
	`
	err := r.db.GetContext(ctx, &review, query, bookingID)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	query := `
		UPDATE reviews 
		SET rating = :rating, comment = :comment
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, review)
	return err
}

func (r *reviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *reviewRepository) GetProviderRating(ctx context.Context, providerID uuid.UUID) (float64, int, error) {
	var rating struct {
		AvgRating   float64 `db:"avg_rating"`
		TotalReviews int    `db:"total_reviews"`
	}
	query := `
		SELECT COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as total_reviews
		FROM reviews WHERE provider_id = $1
	`
	err := r.db.GetContext(ctx, &rating, query, providerID)
	return rating.AvgRating, rating.TotalReviews, err
}

