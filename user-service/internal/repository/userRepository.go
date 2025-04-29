package repository

import (
	"context"
	"errors"
	"user-service/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetUserById(id string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	Delete(id string) error
}

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
	}
}

func (r *UserRepo) Create(user *entity.User) error {
	filter := bson.M{"email": user.Email}
	if err := r.collection.FindOne(context.Background(), filter).Err(); err == nil {
		return errors.New("user already exists")
	}

	res, err := r.collection.InsertOne(context.Background(), bson.M{
		"email":    user.Email,
		"password": user.Password,
	})
	if err != nil {
		return err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserRepo) GetUserById(id string) (*entity.User, error) {
	var user entity.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetAll() ([]*entity.User, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*entity.User
	for cursor.Next(context.Background()) {
		var user entity.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	res, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}
