package main

import (
	"context"
	"encoding/json"

	"qasynda/shared/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewRabbitMQProducer(url string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"chat_messages",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQProducer{
		conn: conn,
		ch:   ch,
		q:    q,
	}, nil
}

func (p *RabbitMQProducer) PublishMessage(msg *Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		context.Background(),
		"",
		p.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (p *RabbitMQProducer) Close() {
	p.ch.Close()
	p.conn.Close()
}

func StartConsumer(ctx context.Context, url string, store *Store) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("failed to connect to rabbitmq consumer", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("failed to open channel", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"chat_messages",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("failed to declare queue", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("failed to register consumer", err)
		return
	}

	logger.Info("Chat Consumer started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Chat Consumer stopping...")
			return
		case d, ok := <-msgs:
			if !ok {
				logger.Info("Chat Consumer channel closed")
				return
			}
			var msg Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				logger.Error("failed to unmarshal message", err)
				continue
			}

			if err := store.SaveMessage(context.Background(), &msg); err != nil {
				logger.Error("failed to save message to db", err)
			}
		}
	}
}
