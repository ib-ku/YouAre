package main

import (
	"log"

	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	userpb "user-service/pkg/gen/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to user-service: %v", err)
	}
	defer conn.Close()

	userClient := userpb.NewUserServiceClient(conn)

	r := gin.Default()
	userHandler := handler.NewUserHandler(userClient)

	// public routes
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// protected routes
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/profile/:id", userHandler.GetProfile)
		auth.DELETE("/profile/:id", userHandler.DeleteUser)
	}

	log.Println("API Gateway running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run HTTP server: %v", err)
	}
}
