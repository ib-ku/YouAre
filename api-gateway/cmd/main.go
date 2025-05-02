package main

import (
	"log"

	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	orderpb "order_service/pkg/gen/order"
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

	// order-service (new connection)
	orderConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure()) // Assuming order service runs on port 50053
	if err != nil {
		log.Fatalf("could not connect to order-service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn) // Order service client

	// HTTP routes
	r := gin.Default()

	// Handlers for user, product, and order services
	productHandler := handler.NewProductHandler(productClient)
	userHandler := handler.NewUserHandler(userClient)
	orderHandler := handler.NewOrderHandler(orderClient) // New order handler

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

		// New routes for orders
		auth.GET("/orders/:user_id", orderHandler.GetOrdersByUser) // Get orders by user ID
		auth.POST("/orders", orderHandler.CreateOrder)             // Create a new order
		auth.DELETE("/orders/:order_id", orderHandler.DeleteOrder) // Delete an order
	}

	log.Println("API Gateway running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run HTTP server: %v", err)
	}
}
