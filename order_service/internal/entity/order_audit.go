package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderAudit struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OrderID   primitive.ObjectID `bson:"order_id"`
	Action    string             `bson:"action"`
	Timestamp time.Time          `bson:"timestamp"`
}
