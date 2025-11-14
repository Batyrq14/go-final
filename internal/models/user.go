package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleClient   UserRole = "client"
	RoleProvider UserRole = "provider"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         UserRole  `db:"role" json:"role"`
	FullName     string    `db:"full_name" json:"full_name"`
	Phone        string    `db:"phone" json:"phone"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type RegisterRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role" validate:"required,oneof=client provider"`
	FullName string   `json:"full_name" validate:"required"`
	Phone    string   `json:"phone" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string    `json:"token"`
	User  *User     `json:"user"`
}

