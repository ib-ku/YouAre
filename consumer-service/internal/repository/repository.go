package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("orders"),
	}
}

type OrderRepository interface {
	UpdateTotalPrice(orderID string, totalPrice float64) error
}

func (r *MongoRepository) UpdateTotalPrice(orderID string, totalPrice float64) error {
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"total_price": totalPrice}}

	_, err = r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}
