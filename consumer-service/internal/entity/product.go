package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `json:"name"`
	Price float64            `json:"price"`
	Stock int                `json:"stock"`
}
