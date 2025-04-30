package client

import (
	"log"

	userpb "user-service/pkg/gen/user"

	"google.golang.org/grpc"
)

func NewUserServiceClient(addr string) userpb.UserServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to user-service: %v", err)
	}

	return userpb.NewUserServiceClient(conn)
}
