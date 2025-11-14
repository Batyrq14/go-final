package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	IconURL     string    `db:"icon_url" json:"icon_url"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateServiceRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
}

