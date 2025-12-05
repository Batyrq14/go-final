package main

import (
	"context"

	"github.com/jmoiron/sqlx"
)

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
