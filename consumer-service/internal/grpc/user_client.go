package grpc

import (
	"context"
	"time"

	"consumer-service/internal/logger"
	userpb "user-service/pkg/gen/user"
)

type UserClient struct {
	client userpb.UserServiceClient
	logger *logger.Logger
}

func NewUserClient(client userpb.UserServiceClient) *UserClient {
	return &UserClient{
		client: client,
		logger: logger.GetLogger(),
	}
}

func (u *UserClient) GetProfile(req *userpb.ProfileRequest) (*userpb.UserResponse, error) {
	u.logger.Debug("Making gRPC call to user service for user ID: %s", req.UserId)
	startTime := time.Now()

	resp, err := u.client.GetProfile(context.TODO(), req)
	if err != nil {
		u.logger.Error("gRPC call to user service failed for user ID: %s: %v", req.UserId, err)
		return nil, err
	}

	u.logger.Debug("User service response received in %v for user ID: %s",
		time.Since(startTime), req.UserId)
	return resp, nil
}
