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
	StartSession() (mongo.Session, error)
	CreateOrderWithSession(ctx mongo.SessionContext, order *entity.Order) (*entity.Order, error)
	InsertOrderAudit(ctx mongo.SessionContext, audit entity.OrderAudit) error
}

type mongoRepo struct {
	collection      *mongo.Collection
	auditCollection *mongo.Collection
}

func NewMongoRepo(db *mongo.Database) *mongoRepo {
	return &mongoRepo{
		collection:      db.Collection("orders"),
		auditCollection: db.Collection("audit"),
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

// audit collection
func (r *mongoRepo) CreateOrderWithSession(sessCtx mongo.SessionContext, order *entity.Order) (*entity.Order, error) {
	res, err := r.collection.InsertOne(sessCtx, order)
	if err != nil {
		return nil, err
	}
	order.ID = res.InsertedID.(primitive.ObjectID)
	return order, nil
}

func (r *mongoRepo) InsertOrderAudit(sessCtx mongo.SessionContext, audit entity.OrderAudit) error {
	_, err := r.auditCollection.InsertOne(sessCtx, audit)
	return err
}

func (r *mongoRepo) StartSession() (mongo.Session, error) {
	return r.collection.Database().Client().StartSession()
}
