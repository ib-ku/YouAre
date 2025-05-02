package main

import (
	"YouAre/order_service/infrastructure/database"
	"YouAre/order_service/infrastructure/rabbitmq"
	"YouAre/order_service/transport/grpc"
	"YouAre/order_service/usecase"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare queue
	queueName := "order.created"
	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Setup layers
	orderRepo := db.NewInMemoryDB()
	publisher := rabbitmq.NewRabbitMQPublisher(ch, queueName)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, publisher)
	grpcServer := grpc.NewGRPCServer(orderUsecase)

	// Simulate gRPC request (testing without proto)
	grpcServer.CreateOrder(nil, "product-123", 5)
}
