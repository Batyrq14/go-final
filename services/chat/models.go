package main

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `db:"id" json:"id"`
	SenderID   uuid.UUID `db:"sender_id" json:"sender_id"`
	ReceiverID uuid.UUID `db:"receiver_id" json:"receiver_id"`
	Content    string    `db:"content" json:"content"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
