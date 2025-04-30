package usecase

import (
	"YouAre/order_service/domain"
)

type OrderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUsecase(repo domain.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}

func (u *OrderUsecase) CreateOrder(order *domain.Order) error {
	return u.repo.CreateOrder(order)
}
