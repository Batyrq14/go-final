package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID         uuid.UUID `db:"id" json:"id"`
	BookingID  uuid.UUID `db:"booking_id" json:"booking_id"`
	ClientID   uuid.UUID `db:"client_id" json:"client_id"`
	ProviderID uuid.UUID `db:"provider_id" json:"provider_id"`
	Rating     int       `db:"rating" json:"rating"` // 1-5
	Comment    string    `db:"comment" json:"comment"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	Client     *User     `json:"client,omitempty"`
}

type CreateReviewRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Comment   string    `json:"comment" validate:"required"`
}

