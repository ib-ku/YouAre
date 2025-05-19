package migrations

import "go.mongodb.org/mongo-driver/mongo"

func RunAllMigrations(db *mongo.Database) error {
	if err := AddIsActiveField(db); err != nil {
		return err
	}
	return nil
}
