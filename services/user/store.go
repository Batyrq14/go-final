package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, full_name, phone, created_at, updated_at)
		VALUES (:id, :email, :password_hash, :role, :full_name, :phone, :created_at, :updated_at)
	`
	_, err := s.db.NamedExecContext(ctx, query, user)
	return err
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE email = $1`
	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE id = $1`
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) ListProviders(limit, offset int) ([]*DetailedProvider, error) {
	var providers []*DetailedProvider
	query := `
		SELECT u.id, sp.id as service_provider_id, u.email, u.full_name, u.role, u.phone, 
		       COALESCE(sp.hourly_rate, 0) as hourly_rate,
		       COALESCE(sp.experience_years, 0) as experience_years,
		       COALESCE(sp.location, '') as location,
		       COALESCE(sp.bio, '') as bio,
		       COALESCE(sp.is_available, false) as is_available,
		       COALESCE(sp.rating, 0) as rating
		FROM users u
		LEFT JOIN service_providers sp ON u.id = sp.user_id
		WHERE u.role = 'provider'
		LIMIT $1 OFFSET $2
	`
	err := s.db.Select(&providers, query, limit, offset)
	return providers, err
}
