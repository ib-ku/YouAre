package entity

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `json:"name"`
	Price float64            `json:"price"`
	Stock int                `json:"stock"`
}

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInternalServerError = errors.New("internal server error")
)
