package main

import (
	"context"
	"log"
	"net"
	"product-service/internal/repository"
	"product-service/internal/usecase"

	productgrpc "product-service/internal/delivery/grpc"
	productpb "product-service/pkg/gen/product"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port     = ":50052"
	mongoURI = "mongodb://localhost:27017"
	dbName   = "YouAre"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database(dbName)
	productRepo := repository.NewProductRepo(db)

	productUC := usecase.NewProductUseCase(productRepo)

	grpcServer := grpc.NewServer()
	productHandler := productgrpc.NewProductHandler(productUC)

	productpb.RegisterProductServiceServer(grpcServer, productHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}
	log.Printf("Product gRPC service started on %s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
