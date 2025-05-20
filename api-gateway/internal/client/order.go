package client

import (
	orderpb "order_service/pkg/gen/order" // Import the generated order service protobuf

	"log"

	"google.golang.org/grpc"
)

// NewOrderClient creates a new GRPC client connection to the order service
func NewOrderClient(addr string) orderpb.OrderServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to order-service: %v", err)
	}
	// conn.Close() should be called in main.go or when you're done using the client
	return orderpb.NewOrderServiceClient(conn)
}
