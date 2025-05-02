package db

import (
	"YouAre/order_service/domain"
	"fmt"
)

type InMemoryDB struct {
	orders map[string]*domain.Order
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		orders: make(map[string]*domain.Order),
	}
}

func (db *InMemoryDB) CreateOrder(order *domain.Order) error {
	db.orders[order.ID] = order
	fmt.Printf("Order saved: %+v\n", order)
	return nil
}
