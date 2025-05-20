package rabbitmq

import (
	"encoding/json"
	"order_service/internal/entity"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewProducer(conn *amqp091.Connection) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare("order.created", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Producer) PublishOrderCreated(order *entity.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
