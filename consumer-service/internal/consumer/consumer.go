package consumer

import (
	"consumer-service/internal/entity"
	"consumer-service/internal/grpc"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn          *amqp091.Connection
	productClient *grpc.ProductClient
}

func NewConsumer(conn *amqp091.Connection, productClient *grpc.ProductClient) *Consumer {
	return &Consumer{
		conn:          conn,
		productClient: productClient,
	}
}

func (c *Consumer) StartConsuming() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		"order.created", // имя очереди
		true,            // durable
		false,           // autoDelete
		false,           // exclusive
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		"order.created", // очередь
		"",              // consumer
		true,            // auto-ack
		false, false, false, nil,
	)
	if err != nil {
		return err
	}

	log.Println(" [*] Waiting for messages from order.created...")

	for msg := range msgs {
		var order entity.Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Printf("Failed to parse order message: %v", err)
			continue
		}

		log.Printf("Received order: %+v", order)

		err = c.productClient.DecreaseStock(order.ProductID, int32(order.Quantity))
		if err != nil {
			log.Printf("Failed to decrease stock: %v", err)
		} else {
			log.Printf("Successfully decreased stock for product %s", order.ProductID)
		}
	}

	return nil
}
