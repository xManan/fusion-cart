package services

import (
	"context"
	"time"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"github.com/xManan/fusion-cart/internal/app/db"
	"github.com/xManan/fusion-cart/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindItemByRef(itemRef string) (models.Item, error) {
	var item models.Item
	err := db.DB.Collection(item.Collection()).FindOne(context.TODO(), bson.M{"itemRef": itemRef}).Decode(&item)
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

func ItemExists(id primitive.ObjectID) error {
	cursor := db.DB.Collection(constants.ItemsCollection).FindOne(context.TODO(), bson.M{"_id": id})
	if cursor.Err() == mongo.ErrNoDocuments {
		return constants.ErrItemNotFound
	}
	return cursor.Err()
}

func FindItemInBasket(itemId, userId, basketId primitive.ObjectID) error {
	filter := bson.D{
		{Key: "_id", Value: userId},
		{Key: "baskets", Value: bson.M{
			"$elemMatch": bson.D{
				{Key: "_id", Value: basketId},
				{Key: "items", Value: bson.M{
					"$elemMatch": bson.M{
						"itemId": itemId,
					},
				}},
			},
		}},
	}
	cursor := db.DB.Collection(constants.UsersCollection).FindOne(context.TODO(), filter)
	if cursor.Err() == mongo.ErrNoDocuments {
		return constants.ErrItemNotFoundInBasket
	}
	return cursor.Err()
}
