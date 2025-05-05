package consumer

import (
	"consumer-service/internal/entity"
	"consumer-service/internal/grpc"
	"consumer-service/internal/repository"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn          *amqp091.Connection
	productClient *grpc.ProductClient
	repo          repository.OrderRepository
}

func NewConsumer(conn *amqp091.Connection, productClient *grpc.ProductClient, repo repository.OrderRepository) *Consumer {
	return &Consumer{
		conn:          conn,
		productClient: productClient,
		repo:          repo,
	}
}

var (
	lastLogTime time.Time
	logInterval = 10 * time.Second
)

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
		false,           // auto-ack
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
			msg.Nack(false, false)
			continue
		}

		product, err := c.productClient.GetProduct(order.ProductID)
		if err != nil {
			if time.Since(lastLogTime) > logInterval {
				log.Printf("Failed to get product: %v", err)
				lastLogTime = time.Now()
			}
			msg.Nack(false, true)
			continue
		}

		order.TotalPrice = product.Product.Price * float64(order.Quantity)

		err = c.repo.UpdateTotalPrice(order.ID.Hex(), float64(order.TotalPrice))
		if err != nil {
			log.Printf("Failed to update total price in DB: %v", err)
			msg.Nack(false, false)
			continue
		}

		err = c.productClient.DecreaseStock(order.ProductID, int32(order.Quantity))
		if err != nil {
			if time.Since(lastLogTime) > logInterval {
				log.Printf("Failed to decrease stock: %v", err)
				lastLogTime = time.Now()
			}
			msg.Nack(false, true)
			continue
		}

		log.Printf("Successfully processed order. ProductID: %s, TotalPrice: %.2f", order.ProductID, order.TotalPrice)
		msg.Ack(false)

	}

	return nil
}
