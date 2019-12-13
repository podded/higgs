package higgs

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Database *mongo.Client
	DBName string
}

func GetDatabaseHandle(config Configuration) (db *DB, err error) {

	clientOptions := options.Client().ApplyURI(config.Database.URI)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return
	}
	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		return
	}

	return &DB{Database: client, DBName:config.Database.Database}, nil
}