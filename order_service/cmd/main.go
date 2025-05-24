package main

import (
	"context"
	"log"
	"net"

	ordergrpc "order_service/internal/delivery/grpc"
	"order_service/internal/rabbitmq"
	"order_service/internal/repository"
	"order_service/internal/usecase"
	orderpb "order_service/pkg/gen/order"
	productpb "product-service/pkg/gen/product"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "YouAre"
)

func main() {
	// Connect to ProductService
	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to ProductService: %v", err)
	}
	defer productConn.Close()

	// Create ProductService client
	productClient := productpb.NewProductServiceClient(productConn)

	// Initialize repository and use case
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	db := client.Database(dbName)

	repo := repository.NewMongoRepo(db)

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	producer, err := rabbitmq.NewProducer(conn)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ producer: %v", err)
	}

	uc := usecase.NewOrderUseCase(repo, productClient, producer)

	// Start TCP listener
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed to listen on port 5000: %v", err)
	}

	// Create gRPC server
	server := grpc.NewServer()

	// Register gRPC service handler
	orderHandler := ordergrpc.NewHandler(uc) // Make sure this exists!
	orderpb.RegisterOrderServiceServer(server, orderHandler)

	// Register reflection
	reflection.Register(server)

	log.Println("OrderService is running on port :5000")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
