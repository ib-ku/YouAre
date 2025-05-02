package repository

import (
	"errors"
	"order_service/internal/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRepository interface {
	CreateOrder(order *entity.Order) (*entity.Order, error)
	GetOrder(id string) (*entity.Order, error)
	GetAllOrders() ([]*entity.Order, error)
	UpdateOrder(id string, quantity int32) (*entity.Order, error)
	DeleteOrder(id string) error
}

type memoryRepo struct {
	orders map[string]*entity.Order
}

func NewMemoryRepo() OrderRepository {
	return &memoryRepo{
		orders: make(map[string]*entity.Order),
	}
}

func (r *memoryRepo) CreateOrder(order *entity.Order) (*entity.Order, error) {
	order.ID = primitive.NewObjectID()
	r.orders[order.ID.Hex()] = order
	return order, nil
}

func (r *memoryRepo) GetOrder(id string) (*entity.Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (r *memoryRepo) GetAllOrders() ([]*entity.Order, error) {
	var result []*entity.Order
	for _, order := range r.orders {
		result = append(result, order)
	}
	return result, nil
}

func (r *memoryRepo) UpdateOrder(id string, quantity int32) (*entity.Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	order.Quantity = int(quantity)
	return order, nil
}

func (r *memoryRepo) DeleteOrder(id string) error {
	if _, ok := r.orders[id]; !ok {
		return errors.New("order not found")
	}
	delete(r.orders, id)
	return nil
}
