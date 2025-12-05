package main

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	IconURL     string    `db:"icon_url"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Booking struct {
	ID            uuid.UUID `db:"id"`
	ClientID      uuid.UUID `db:"client_id"`
	ProviderID    uuid.UUID `db:"provider_id"`
	ServiceID     uuid.UUID `db:"service_id"`
	ScheduledDate time.Time `db:"scheduled_date"`
	DurationHours float64   `db:"duration_hours"`
	Status        string    `db:"status"`
	TotalPrice    float64   `db:"total_price"`
	Notes         string    `db:"notes"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
