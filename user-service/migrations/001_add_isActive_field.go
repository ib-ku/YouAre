package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddIsActiveField(db *mongo.Database) error {
	_, err := db.Collection("users").UpdateMany(
		context.TODO(),
		bson.M{},
		bson.M{"$set": bson.M{"isActive": true}},
	)
	return err
}
