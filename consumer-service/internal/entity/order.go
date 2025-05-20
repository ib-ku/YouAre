package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	ProductID  string             `bson:"product_id"`
	Quantity   int                `bson:"quantity"`
	TotalPrice float64            `bson:"total_price"`
	CreatedAt  time.Time          `bson:"created_at"`
}
