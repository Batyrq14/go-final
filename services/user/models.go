package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	FullName     string    `db:"full_name"`
	Phone        string    `db:"phone"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type DetailedProvider struct {
	ID                uuid.UUID `db:"id"`
	ServiceProviderID uuid.UUID `db:"service_provider_id"`
	Email             string    `db:"email"`
	FullName          string    `db:"full_name"`
	Role              string    `db:"role"`
	Phone             string    `db:"phone"`
	HourlyRate        float64   `db:"hourly_rate"`
	ExperienceYears   int32     `db:"experience_years"`
	Location          string    `db:"location"`
	Bio               string    `db:"bio"`
	IsAvailable       bool      `db:"is_available"`
	Rating            float64   `db:"rating"`
}
