package rabbitmq

import (
	"YouAre/order_service/domain"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type EventPublisher interface {
	PublishOrderCreated(order *domain.Order) error
}

type RabbitMQPublisher struct {
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQPublisher(channel *amqp.Channel, queueName string) *RabbitMQPublisher {
	return &RabbitMQPublisher{channel: channel, queueName: queueName}
}

func (p *RabbitMQPublisher) PublishOrderCreated(order *domain.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",          // default exchange
		p.queueName, // routing key = queue name
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("[RabbitMQ] Order Published: %s at %v", order.ID, order.ID)
	return nil
}
