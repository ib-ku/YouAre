package repository

import (
	"context"
	"errors"
	"product-service/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	GetProductByID(id string) (*entity.Product, error)
	GetAllProducts() ([]*entity.Product, error)
	Update(product *entity.Product) (*entity.Product, error)
	Decrease(id string, quantity int) (*entity.Product, error)
	Delete(id string) error
}

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{
		collection: db.Collection("products"),
	}
}

// realization
func (r *ProductRepo) Create(product *entity.Product) error {
	filter := bson.M{"id": product.ID}
	if err := r.collection.FindOne(context.Background(), filter).Err(); err == nil {
		return errors.New("product already exists")
	}

	res, err := r.collection.InsertOne(context.Background(), bson.M{
		"name":  product.Name,
		"price": product.Price,
		"stock": product.Stock,
	})
	if err != nil {
		return err
	}
	product.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ProductRepo) GetProductByID(id string) (*entity.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product entity.Product
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepo) GetAllProducts() ([]*entity.Product, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var products []*entity.Product
	for cursor.Next(context.Background()) {
		var product entity.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *ProductRepo) Update(product *entity.Product) (*entity.Product, error) {
	filter := bson.M{"_id": product.ID}
	update := bson.M{
		"$set": bson.M{
			"name":  product.Name,
			"price": product.Price,
			"stock": product.Stock,
		},
	}

	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepo) Decrease(id string, quantity int) (*entity.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "stock": bson.M{"$gte": quantity}}
	update := bson.M{"$inc": bson.M{"stock": -quantity}}

	result := r.collection.FindOneAndUpdate(
		context.Background(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updated entity.Product
	if err := result.Decode(&updated); err != nil {
		return nil, errors.New("not enough stock or product not found")
	}

	return &updated, nil
}

func (r *ProductRepo) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}
