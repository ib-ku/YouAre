package usecase

import (
	"order_service/internal/entity"
	"order_service/internal/repository"
	productpb "product-service/pkg/gen/product"
)

type OrderUseCase struct {
	repo          repository.OrderRepository
	productClient productpb.ProductServiceClient
}

func NewOrderUseCase(repo repository.OrderRepository, productClient productpb.ProductServiceClient) *OrderUseCase {
	return &OrderUseCase{
		repo:          repo,
		productClient: productClient,
	}
}

func (u *OrderUseCase) CreateOrder(order *entity.Order) (*entity.Order, error) {
	return u.repo.CreateOrder(order)
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
