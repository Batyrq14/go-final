package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceProvider struct {
	ID             uuid.UUID `db:"id" json:"id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	Bio            string    `db:"bio" json:"bio"`
	HourlyRate     float64   `db:"hourly_rate" json:"hourly_rate"`
	ExperienceYears int      `db:"experience_years" json:"experience_years"`
	Location       string    `db:"location" json:"location"`
	IsAvailable    bool      `db:"is_available" json:"is_available"`
	Rating         float64   `db:"rating" json:"rating"`
	TotalReviews   int       `db:"total_reviews" json:"total_reviews"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	User           *User     `json:"user,omitempty"`
	Services       []*Service `json:"services,omitempty"`
}

type UpdateProviderRequest struct {
	Bio            *string  `json:"bio"`
	HourlyRate     *float64 `json:"hourly_rate" validate:"omitempty,min=0"`
	ExperienceYears *int    `json:"experience_years" validate:"omitempty,min=0"`
	Location       *string  `json:"location"`
	IsAvailable    *bool    `json:"is_available"`
}

type ProviderFilter struct {
	ServiceID  *uuid.UUID
	City       *string
	MinRating  *float64
	IsAvailable *bool
	Page       int
	Limit      int
	SortBy     string // "rating", "experience", "hourly_rate"
	SortOrder  string // "asc", "desc"
}

type ProviderListResponse struct {
	Data       []*ServiceProvider `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

