package models

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusAccepted  BookingStatus = "accepted"
	StatusRejected  BookingStatus = "rejected"
	StatusCompleted BookingStatus = "completed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID            uuid.UUID        `db:"id" json:"id"`
	ClientID      uuid.UUID        `db:"client_id" json:"client_id"`
	ProviderID    uuid.UUID        `db:"provider_id" json:"provider_id"`
	ServiceID     uuid.UUID        `db:"service_id" json:"service_id"`
	ScheduledDate time.Time        `db:"scheduled_date" json:"scheduled_date"`
	DurationHours float64          `db:"duration_hours" json:"duration_hours"`
	Status        BookingStatus    `db:"status" json:"status"`
	TotalPrice    float64          `db:"total_price" json:"total_price"`
	Notes         string           `db:"notes" json:"notes"`
	CreatedAt     time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time        `db:"updated_at" json:"updated_at"`
	Client        *User            `json:"client,omitempty"`
	Provider      *ServiceProvider `json:"provider,omitempty"`
	Service       *Service         `json:"service,omitempty"`
}

type CreateBookingRequest struct {
	ProviderID    uuid.UUID `json:"provider_id" validate:"required"`
	ServiceID     uuid.UUID `json:"service_id" validate:"required"`
	ScheduledDate string    `json:"scheduled_date" validate:"required"`
	DurationHours float64   `json:"duration_hours" validate:"required,min=0.5"`
	Notes         string    `json:"notes"`
}

type UpdateBookingStatusRequest struct {
	Status BookingStatus `json:"status" validate:"required,oneof=accepted rejected completed cancelled"`
}
