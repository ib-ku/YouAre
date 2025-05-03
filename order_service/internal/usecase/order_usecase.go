package usecase

import (
	"context"
	"log"
	"order_service/internal/entity"
	"order_service/internal/rabbitmq"
	"order_service/internal/repository"
	productpb "product-service/pkg/gen/product"
	"time"
)

type OrderUseCase struct {
	repo          repository.OrderRepository
	productClient productpb.ProductServiceClient
	producer      *rabbitmq.Producer
}

func NewOrderUseCase(repo repository.OrderRepository, productClient productpb.ProductServiceClient, producer *rabbitmq.Producer) *OrderUseCase {
	return &OrderUseCase{
		repo:          repo,
		productClient: productClient,
		producer:      producer,
	}
}

func (u *OrderUseCase) CreateOrder(order *entity.Order) (*entity.Order, error) {
	// Получаем информацию о продукте через gRPC
	resp, err := u.productClient.Get(context.TODO(), &productpb.ProductRequest{
		Id: order.ProductID,
	})
	if err != nil {
		return nil, err
	}

	order.TotalPrice = float64(order.Quantity) * resp.Product.Price
	order.CreatedAt = time.Now()

	createdOrder, err := u.repo.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	err = u.producer.PublishOrderCreated(createdOrder)
	if err != nil {
		log.Printf("Failed to publish order.created event: %v", err)
	}

	return createdOrder, nil
}

func (u *OrderUseCase) GetOrder(id string) (*entity.Order, error) {
	return u.repo.GetOrder(id)
}

func (u *OrderUseCase) GetAllOrders() ([]*entity.Order, error) {
	return u.repo.GetAllOrders()
}

func (u *OrderUseCase) UpdateOrder(id string, quantity int32) (*entity.Order, error) {
	return u.repo.UpdateOrder(id, quantity)
}

func (u *OrderUseCase) DeleteOrder(id string) error {
	return u.repo.DeleteOrder(id)
}
