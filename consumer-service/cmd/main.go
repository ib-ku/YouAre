package main

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"consumer-service/internal/consumer"
	consumerGrpc "consumer-service/internal/grpc"
	"consumer-service/internal/repository"
	productpb "product-service/pkg/gen/product"
)

var mongoURI = "mongodb://localhost:27017"
var dbName = "YouAre"

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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	db := client.Database(dbName)
	repo := repository.NewMongoRepository(db)

	productClient := consumerGrpc.NewProductClient(productpb.NewProductServiceClient(grpcConn))

	consumer := consumer.NewConsumer(conn, productClient, repo)
	if err := consumer.StartConsuming(); err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
