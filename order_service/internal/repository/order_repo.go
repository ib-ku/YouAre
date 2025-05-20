package repository

import (
	"context"
	"order_service/internal/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	CreateOrder(order *entity.Order) (*entity.Order, error)
	GetOrder(id string) (*entity.Order, error)
	GetAllOrders() ([]*entity.Order, error)
	UpdateOrder(id string, quantity int32) (*entity.Order, error)
	DeleteOrder(id string) error
}

type mongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(db *mongo.Database) *mongoRepo {
	return &mongoRepo{
		collection: db.Collection("orders"),
	}
}

func (r *mongoRepo) CreateOrder(order *entity.Order) (*entity.Order, error) {
	order.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(context.TODO(), order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *mongoRepo) GetOrder(id string) (*entity.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order entity.Order
	err = r.collection.FindOne(context.TODO(), primitive.M{"_id": objID}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *mongoRepo) GetAllOrders() ([]*entity.Order, error) {
	cursor, err := r.collection.Find(context.TODO(), primitive.M{})
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

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *mongoRepo) UpdateOrder(id string, quantity int32) (*entity.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := primitive.M{
		"$set": primitive.M{"quantity": quantity},
	}

	_, err = r.collection.UpdateByID(context.TODO(), objID, update)
	if err != nil {
		return nil, err
	}

	// Возвращаем обновлённый заказ
	return r.GetOrder(id)
}

func (r *mongoRepo) DeleteOrder(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.TODO(), primitive.M{"_id": objID})
	return err
}
