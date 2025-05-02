package main

import (
	"log"
	"net"

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
	// Connect to ProductService
	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to ProductService: %v", err)
	}
	defer productConn.Close()

	// Create ProductService client
	productClient := productpb.NewProductServiceClient(productConn)

	// Initialize repository and use case
	repo := repository.NewMemoryRepo()
	uc := usecase.NewOrderUseCase(repo, productClient)

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

	log.Println("OrderService is running on port :5000...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
