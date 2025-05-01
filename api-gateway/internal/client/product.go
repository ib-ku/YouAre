package client

import (
	productpb "product-service/pkg/gen/product"

	"log"

	"google.golang.org/grpc"
)

func NewProductClient(addr string) productpb.ProductServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to product-service: %v", err)
	}
	// conn.Close() вызывается в main
	return productpb.NewProductServiceClient(conn)
}
