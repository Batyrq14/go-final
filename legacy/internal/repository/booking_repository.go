package repository

import (
	"context"
	"qasynda/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*models.Booking, error)
	GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.Booking, error)
	Update(ctx context.Context, booking *models.Booking) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookingStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type bookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	query := `
		INSERT INTO bookings (id, client_id, provider_id, service_id, scheduled_date, duration_hours, status, total_price, notes, created_at, updated_at)
		VALUES (:id, :client_id, :provider_id, :service_id, :scheduled_date, :duration_hours, :status, :total_price, :notes, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, booking)
	return err
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	query := `
		SELECT id, client_id, provider_id, service_id, scheduled_date, duration_hours, status, total_price, notes, created_at, updated_at
		FROM bookings WHERE id = $1
	`
	err := r.db.GetContext(ctx, &booking, query, id)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*models.Booking, error) {
	var bookings []*models.Booking
	query := `
		SELECT id, client_id, provider_id, service_id, scheduled_date, duration_hours, status, total_price, notes, created_at, updated_at
		FROM bookings WHERE client_id = $1 ORDER BY scheduled_date DESC
	`
	err := r.db.SelectContext(ctx, &bookings, query, clientID)
	return bookings, err
}

func (r *bookingRepository) GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*models.Booking, error) {
	var bookings []*models.Booking
	query := `
		SELECT id, client_id, provider_id, service_id, scheduled_date, duration_hours, status, total_price, notes, created_at, updated_at
		FROM bookings WHERE provider_id = $1 ORDER BY scheduled_date DESC
	`
	err := r.db.SelectContext(ctx, &bookings, query, providerID)
	return bookings, err
}

func (r *bookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	query := `
		UPDATE bookings 
		SET scheduled_date = :scheduled_date, duration_hours = :duration_hours, 
		    status = :status, total_price = :total_price, notes = :notes, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, booking)
	return err
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookingStatus) error {
	query := `UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *bookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM bookings WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

