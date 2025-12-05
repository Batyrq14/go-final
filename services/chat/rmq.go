package main

import (
	"context"
	"encoding/json"
	"log"

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
		"chat_messages", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
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
		"",       // exchange
		p.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (p *RabbitMQProducer) Close() {
	p.ch.Close()
	p.conn.Close()
}

// Consumer that writes to DB
func StartConsumer(url string, store *Store) {
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
		"chat_messages", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		logger.Error("failed to declare queue", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.Error("failed to register consumer", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				logger.Error("failed to unmarshal message", err)
				continue
			}

			// Save to DB
			if err := store.SaveMessage(context.Background(), &msg); err != nil {
				logger.Error("failed to save message to db", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
