package usecase

import (
	"YouAre/order_service/domain"
	"YouAre/order_service/infrastructure/rabbitmq"
)

type OrderUsecase struct {
	repo      domain.OrderRepository
	publisher rabbitmq.EventPublisher
}

func NewOrderUsecase(repo domain.OrderRepository, publisher rabbitmq.EventPublisher) *OrderUsecase {
	return &OrderUsecase{repo: repo, publisher: publisher}
}

func (u *OrderUsecase) CreateOrder(order *domain.Order) error {
	if err := u.repo.CreateOrder(order); err != nil {
		return err
	}
	return u.publisher.PublishOrderCreated(order)
}
