package main

import (
	"log"

	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	productpb "product-service/pkg/gen/product"
	userpb "user-service/pkg/gen/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// user-service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to user-service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	// product-service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to product-service: %v", err)
	}
	defer productConn.Close()
	productClient := productpb.NewProductServiceClient(productConn)

	// HTTP routes
	r := gin.Default()

	productHandler := handler.NewProductHandler(productClient)
	userHandler := handler.NewUserHandler(userClient)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/profile/:id", userHandler.GetProfile)
		auth.DELETE("/profile/:id", userHandler.DeleteUser)

		auth.GET("/products", productHandler.ListProducts)
		auth.GET("/products/:id", productHandler.GetProductByID)
		auth.POST("/products", productHandler.CreateProduct)
		auth.DELETE("/products/:id", productHandler.DeleteProduct)

	}
	log.Println("API Gateway running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run HTTP server: %v", err)
	}
}
