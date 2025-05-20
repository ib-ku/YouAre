package main

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"consumer-service/internal/consumer"
	consumerGrpc "consumer-service/internal/grpc"
	"consumer-service/internal/repository"
	productpb "product-service/pkg/gen/product"
	userpb "user-service/pkg/gen/user"
)

const (
	mongoURI    = "mongodb://localhost:27017"
	dbName      = "YouAre"
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
)

func main() {
	// RabbitMQ connection with retry
	var conn *amqp091.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp091.Dial(rabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connection attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// gRPC connections with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Product Service
	productConn, err := grpc.DialContext(ctx, "localhost:50052",
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to product-service: %v", err)
	}
	defer productConn.Close()

	// User Service
	userConn, err := grpc.DialContext(ctx, "localhost:50051",
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to user-service: %v", err)
	}
	defer userConn.Close()

	// MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()

	// Initialize services
	repo := repository.NewMongoRepository(mongoClient.Database(dbName))
	productClient := consumerGrpc.NewProductClient(productpb.NewProductServiceClient(productConn))
	userClient := consumerGrpc.NewUserClient(userpb.NewUserServiceClient(userConn))

	// Start consumer
	consumer := consumer.NewConsumer(conn, productClient, userClient, repo)
	log.Println("Consumer service started successfully")

	if err := consumer.StartConsuming(); err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
