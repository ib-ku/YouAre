package main

import (
	"log"
	"net"
	"os"
	"time"

	"user-service/internal/cache"
	usergrpc "user-service/internal/delivery/grpc"
	"user-service/internal/repository"
	"user-service/internal/usecase"
	userpb "user-service/pkg/gen/user"

	"context"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port      = ":50051"
	mongoURI  = "mongodb://localhost:27017"
	dbName    = "YouAre"
	jwtSecret string
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in environment")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database(dbName)
	userRepo := repository.NewUserRepo(db)
	redisCache := cache.NewRedisCache("localhost:6379")

	accessTokenTTL := 15 * time.Minute

	authUC := usecase.NewAuthUsecase(userRepo, jwtSecret, accessTokenTTL)
	userUC := usecase.NewUserUsecase(userRepo, redisCache)

	grpcServer := grpc.NewServer()
	userHandler := usergrpc.NewUserHandler(userUC, authUC)

	userpb.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}
	log.Printf("User gRPC service started on %s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
