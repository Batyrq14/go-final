package main

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type IStore interface {
	CreateService(ctx context.Context, service *Service) error
	ListServices(ctx context.Context) ([]*Service, error)
	CreateBooking(ctx context.Context, booking *Booking) error
	ListBookings(ctx context.Context, userID string, role string) ([]*Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID string, status string) error
	GetBooking(ctx context.Context, bookingID string) (*Booking, error)
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateService(ctx context.Context, service *Service) error {
	query := `
		INSERT INTO services (id, name, description, icon_url, created_at, updated_at)
		VALUES (:id, :name, :description, :icon_url, :created_at, :updated_at)
	`
	_, err := s.db.NamedExecContext(ctx, query, service)
	return err
}

func (s *Store) ListServices(ctx context.Context) ([]*Service, error) {
	var services []*Service
	query := `SELECT * FROM services ORDER BY name`
	err := s.db.SelectContext(ctx, &services, query)
	return services, err
}

func (s *Store) CreateBooking(ctx context.Context, booking *Booking) error {
	query := `
		INSERT INTO bookings (id, client_id, provider_id, service_id, scheduled_date, duration_hours, status, total_price, notes, created_at, updated_at)
		VALUES (:id, :client_id, :provider_id, :service_id, :scheduled_date, :duration_hours, :status, :total_price, :notes, :created_at, :updated_at)
	`
	_, err := s.db.NamedExecContext(ctx, query, booking)
	return err
}

func (s *Store) ListBookings(ctx context.Context, userID string, role string) ([]*Booking, error) {
	var bookings []*Booking
	var query string
	if role == "provider" {

		query = `
			SELECT b.* 
			FROM bookings b
			INNER JOIN service_providers sp ON b.provider_id = sp.id
			WHERE sp.user_id = $1 
			ORDER BY b.created_at DESC
		`
	} else {
		query = `SELECT * FROM bookings WHERE client_id = $1 ORDER BY created_at DESC`
	}
	err := s.db.SelectContext(ctx, &bookings, query, userID)
	return bookings, err
}

func (s *Store) UpdateBookingStatus(ctx context.Context, bookingID string, status string) error {
	query := `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, status, bookingID)
	return err
}

func (s *Store) GetBooking(ctx context.Context, bookingID string) (*Booking, error) {
	var booking Booking
	query := `SELECT * FROM bookings WHERE id = $1`
	err := s.db.GetContext(ctx, &booking, query, bookingID)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
