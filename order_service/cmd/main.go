package main

import (
	"context"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	ordergrpc "order_service/internal/delivery/grpc"
	"order_service/internal/repository"
	"order_service/internal/usecase"
	orderpb "order_service/pkg/gen/order"
	productpb "product-service/pkg/gen/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// MongoDB setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Modify this if MongoDB is not on localhost
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	// Repository and UseCase Setup
	repo := repository.NewMongoRepo(mongoClient) // Use MongoRepo
	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to ProductService: %v", err)
	}
	defer productConn.Close()

	productClient := productpb.NewProductServiceClient(productConn)
	uc := usecase.NewOrderUseCase(repo, productClient)

	// gRPC Server Setup
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed to listen on port 5000: %v", err)
	}

	server := grpc.NewServer()
	orderHandler := ordergrpc.NewHandler(uc)
	orderpb.RegisterOrderServiceServer(server, orderHandler)
	reflection.Register(server)

	log.Println("OrderService is running on port :5000...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
