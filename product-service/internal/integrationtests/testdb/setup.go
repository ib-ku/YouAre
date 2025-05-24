package testdb

import (
	"context"
	"log"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *mongo.Database
)

func Init() *mongo.Database {
	if db != nil {
		return db
	}

	connStr := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db = client.Database("test_products")
	return db
}

func GetDB() *mongo.Database {
	if db == nil {
		panic("DB not initialized. Call Init() first")
	}
	return db
}

func Cleanup() {
	if db != nil {
		_ = db.Drop(context.Background())
	}
}

func TestMainWrapper(m *testing.M) {
	Init()
	code := m.Run()
	Cleanup()
	os.Exit(code)
}
