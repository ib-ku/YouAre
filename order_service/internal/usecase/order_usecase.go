package usecase

import (
	"context"
	"encoding/json"
	"log"
	"order_service/internal/cache"
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
	cache         *cache.RedisCache
}

func NewOrderUseCase(repo repository.OrderRepository, productClient productpb.ProductServiceClient, producer *rabbitmq.Producer, redisCache *cache.RedisCache) *OrderUseCase {
	return &OrderUseCase{
		repo:          repo,
		productClient: productClient,
		producer:      producer,
		cache:         redisCache,
	}
}

func (u *OrderUseCase) CreateOrder(order *entity.Order) (*entity.Order, error) {
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
	ctx := context.Background()
	cacheKey := "order:" + id

	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		log.Println("[CACHE HIT] order ID:", id)
		var order entity.Order
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			return &order, nil
		}
		log.Println("[CACHE ERROR] Failed to unmarshal cached data for order ID:", id)
	}

	log.Println("[CACHE MISS] order ID:", id)

	order, err := u.repo.GetOrder(id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(order)
	if err == nil {
		_ = u.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return order, nil
}

func (u *OrderUseCase) GetAllOrders() ([]*entity.Order, error) {
	return u.repo.GetAllOrders()
}

func (u *OrderUseCase) UpdateOrder(id string, quantity int32) (*entity.Order, error) {
	updatedOrder, err := u.repo.UpdateOrder(id, quantity)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	cacheKey := "order:" + id
	if err := u.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("failed to invalidate cache for order %s: %v", id, err)
	}

	return updatedOrder, nil
}

func (u *OrderUseCase) DeleteOrder(id string) error {
	err := u.repo.DeleteOrder(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cacheKey := "order:" + id
	if err := u.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("failed to invalidate cache for order %s: %v", id, err)
	}

	return nil
}
