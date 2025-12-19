package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IStore interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	ListProviders(limit, offset int) ([]*DetailedProvider, error)
	UpdateProviderStatus(ctx context.Context, userID uuid.UUID, isAvailable bool) error
	GetProviderStatus(ctx context.Context, userID uuid.UUID) (bool, error)
}

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (id, email, password_hash, role, full_name, phone, created_at, updated_at)
		VALUES (:id, :email, :password_hash, :role, :full_name, :phone, :created_at, :updated_at)
	`
	_, err = tx.NamedExecContext(ctx, query, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	if user.Role == "provider" {

		spQuery := `
			INSERT INTO service_providers (id, user_id, is_available)
			VALUES ($1, $2, true)
		`
		_, err = tx.ExecContext(ctx, spQuery, uuid.New(), user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE email = $1`
	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
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

func (s *UserStore) UpdateProviderStatus(ctx context.Context, userID uuid.UUID, isAvailable bool) error {

	query := `
		INSERT INTO service_providers (id, user_id, is_available)
		VALUES (uuid_generate_v4(), $2, $1)
		ON CONFLICT (user_id) DO UPDATE
		SET is_available = $1
	`
	_, err := s.db.ExecContext(ctx, query, isAvailable, userID)
	return err
}

func (s *UserStore) GetProviderStatus(ctx context.Context, userID uuid.UUID) (bool, error) {
	var isAvailable bool
	query := `SELECT COALESCE(is_available, false) FROM service_providers WHERE user_id = $1`
	err := s.db.GetContext(ctx, &isAvailable, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isAvailable, nil
}
