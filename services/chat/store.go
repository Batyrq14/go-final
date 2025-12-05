package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) SaveMessage(ctx context.Context, msg *Message) error {
	query := `
		INSERT INTO messages (id, sender_id, receiver_id, content, created_at)
		VALUES (:id, :sender_id, :receiver_id, :content, :created_at)
	`
	_, err := s.db.NamedExecContext(ctx, query, msg)
	return err
}

func (s *Store) GetHistory(ctx context.Context, userID1, userID2 uuid.UUID, limit, offset int) ([]*Message, error) {
	var messages []*Message
	query := `
		SELECT * FROM messages 
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	err := s.db.SelectContext(ctx, &messages, query, userID1, userID2, limit, offset)
	return messages, err
}
