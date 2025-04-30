package main

import (
	db "YouAre/order_service/infrastructure/database"
	"YouAre/order_service/transport/grpc"
	"YouAre/order_service/usecase"
)

func main() {
	// Setup layers
	orderRepo := db.NewInMemoryDB()
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	grpcServer := grpc.NewGRPCServer(orderUsecase)

	// Simulate gRPC request (for testing)
	grpcServer.CreateOrder(nil, "product-123", 5)
}
