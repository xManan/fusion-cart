package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func MongoInit(uri string, dbName string) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	DB = client.Database(dbName)
	if err := DB.Client().Ping(context.TODO(), nil); err != nil {
		return err
	}
	return nil
}

func MongoClose() error {
	return DB.Client().Disconnect(context.TODO())
}
