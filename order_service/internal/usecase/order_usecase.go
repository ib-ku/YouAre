package usecase

import (
	"context"
	"log"
	"order_service/internal/entity"
	"order_service/internal/rabbitmq"
	"order_service/internal/repository"
	productpb "product-service/pkg/gen/product"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
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

// transaction
func (u *OrderUseCase) CreateOrder(order *entity.Order) (*entity.Order, error) {
	order.CreatedAt = time.Now()

	session, err := u.repo.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	result, err := session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		createdOrder, err := u.repo.CreateOrderWithSession(sessCtx, order)
		if err != nil {
			return nil, err
		}

		audit := entity.OrderAudit{
			OrderID:   createdOrder.ID,
			Action:    "created",
			Timestamp: time.Now(),
		}

		err = u.repo.InsertOrderAudit(sessCtx, audit)
		if err != nil {
			return nil, err
		}

		return createdOrder, nil
	})
	if err != nil {
		return nil, err
	}

	createdOrder := result.(*entity.Order)
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
