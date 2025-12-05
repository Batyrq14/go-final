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
