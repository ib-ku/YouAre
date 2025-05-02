package repository

import (
	"context"
	"errors"
	"log"
	"order_service/internal/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoRepo struct with MongoDB collection
type mongoRepo struct {
	collection *mongo.Collection
}

// NewMongoRepo creates a new repository connected to MongoDB
func NewMongoRepo(client *mongo.Client) OrderRepository {
	collection := client.Database("order_service").Collection("orders")
	return &mongoRepo{collection: collection}
}

// CreateOrder inserts an order into MongoDB
func (r *mongoRepo) CreateOrder(order *entity.Order) (*entity.Order, error) {
	order.ID = primitive.NewObjectID()
	order.CreatedAt = time.Now()

	log.Printf("Inserting order: %+v", order)
	_, err := r.collection.InsertOne(context.TODO(), order)
	if err != nil {
		return nil, err
	}
	log.Println("Order inserted successfully.")
	return order, nil
}

// GetOrder retrieves an order from MongoDB by ID
func (r *mongoRepo) GetOrder(id string) (*entity.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order entity.Order
	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		return nil, errors.New("order not found")
	}
	return &order, nil
}

// GetAllOrders retrieves all orders from MongoDB
func (r *mongoRepo) GetAllOrders() ([]*entity.Order, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var orders []*entity.Order
	for cursor.Next(context.TODO()) {
		var order entity.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

// UpdateOrder updates an existing order in MongoDB
func (r *mongoRepo) UpdateOrder(id string, quantity int32) (*entity.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": bson.M{"quantity": quantity}}
	_, err = r.collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}

	return r.GetOrder(id)
}

func (r *mongoRepo) DeleteOrder(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("order not found")
	}
	return nil
}

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
