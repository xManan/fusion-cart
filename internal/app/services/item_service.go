package services

import (
	"context"
	"time"

	"github.com/xManan/fusion-cart/internal/app/db"
	"github.com/xManan/fusion-cart/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindItemByRef(itemRef string) (models.Item, error) {
	var item models.Item
	err := db.DB.Collection(item.Collection()).FindOne(context.TODO(), bson.M{ "itemRef": itemRef }).Decode(&item)
	if err != nil {
		return models.Item{}, err
	}
	return item, nil
}

func SaveItem(item *models.Item) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	res, err := db.DB.Collection(item.Collection()).InsertOne(context.TODO(), *item)
	if err != nil {
		return err
	}
	oid := res.InsertedID.(primitive.ObjectID)
	item.Id = oid
	return nil
}
