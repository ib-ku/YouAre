package grpc

import (
	"YouAre/order_service/domain"
	"YouAre/order_service/usecase"
	"context"
	"log"

	"github.com/google/uuid"
)

type GRPCServer struct {
	orderUsecase *usecase.OrderUsecase
}

func NewGRPCServer(orderUsecase *usecase.OrderUsecase) *GRPCServer {
	return &GRPCServer{orderUsecase: orderUsecase}
}

func (s *GRPCServer) CreateOrder(ctx context.Context, productID string, quantity int) error {
	order := &domain.Order{
		ID:        uuid.New().String(),
		ProductID: productID,
		Quantity:  quantity,
	}

	if err := s.orderUsecase.CreateOrder(order); err != nil {
		log.Printf("Failed to create order: %v", err)
		return err
	}

	log.Printf("Order created successfully: %+v", order)
	return nil
}
