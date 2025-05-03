package main

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"

	"consumer-service/internal/consumer"
	consumerGrpc "consumer-service/internal/grpc"
	productpb "product-service/pkg/gen/product"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %v", err)
	}
	defer conn.Close()

	grpcConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to product-service: %v", err)
	}
	defer grpcConn.Close()

	productClient := consumerGrpc.NewProductClient(productpb.NewProductServiceClient(grpcConn))

	consumer := consumer.NewConsumer(conn, productClient)
	if err := consumer.StartConsuming(); err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
