package main

import (
	"context"
	"log"
	"net"

	"order_service/internal/cache"
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
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	db := client.Database(dbName)
	repo := repository.NewMongoRepo(db)

	// Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %v", err)
	}
	producer, err := rabbitmq.NewProducer(conn)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ producer: %v", err)
	}

	// Connect to ProductService
	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to ProductService: %v", err)
	}
	defer productConn.Close()
	productClient := productpb.NewProductServiceClient(productConn)

	// Initialize Redis cache
	redisCache := cache.NewRedisCache("localhost:6379")

	// Initialize use case with all dependencies
	uc := usecase.NewOrderUseCase(repo, productClient, producer, redisCache)

	// Start TCP listener
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Failed to listen on port 5000: %v", err)
	}

	// Create gRPC server
	server := grpc.NewServer()

	// Register gRPC handler
	orderHandler := ordergrpc.NewHandler(uc)
	orderpb.RegisterOrderServiceServer(server, orderHandler)

	// Enable reflection
	reflection.Register(server)

	log.Println("OrderService is running on port :5000...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
