package services

import (
	"context"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"github.com/xManan/fusion-cart/internal/app/db"
	"github.com/xManan/fusion-cart/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddBasketToUserBaskets(userId primitive.ObjectID, name string) error {
	basket := models.Basket{
		Id:    primitive.NewObjectID(),
		Name:  name,
		Items: []models.BasketItem{},
	}
	update := bson.M{"$push": bson.M{"baskets": basket}}
	_, err := db.DB.Collection(constants.UsersCollection).UpdateByID(context.TODO(), userId, update)
	return err
}

func GetBaskets(userId primitive.ObjectID) ([]models.Basket, error) {
	var user models.User
	opts := options.FindOne().SetProjection(bson.M{"baskets": 1})
	err := db.DB.Collection(constants.UsersCollection).FindOne(context.TODO(), bson.M{"_id": userId}, opts).Decode(&user)
	return user.Baskets, err

}

func GetBasket(userId, basketId primitive.ObjectID) (models.Basket, error) {
	var user models.User
	filter := bson.D{
		{
			Key:   "_id",
			Value: userId,
		},
		{
			Key: "baskets",
			Value: bson.M{
				"$elemMatch": bson.M{
					"_id": basketId,
				},
			},
		},
	}

	opts := options.FindOne().SetProjection(bson.M{"baskets": 1})
	err := db.DB.Collection(constants.UsersCollection).FindOne(context.TODO(), filter, opts).Decode(&user)
	return user.Baskets[0], err
}

func UpdateBasket(userId, basketId primitive.ObjectID, name string) error {
	filter := bson.M{"_id": userId, "baskets": bson.M{"$elemMatch": bson.M{"_id": basketId}}}
	update := bson.M{"$set": bson.M{"baskets.$.name": name}}
	res, err := db.DB.Collection(constants.UsersCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return constants.ErrBasketNotFound
	}
	return nil
}

func DeleteBasket(userId, basketId primitive.ObjectID) error {
	filter := bson.M{"_id": userId}
	update := bson.M{
		"$pull": bson.M{
			"baskets": bson.M{
				"_id": basketId,
			},
		},
	}
	res, err := db.DB.Collection(constants.UsersCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return constants.ErrBasketNotFound
	}
	return nil
}

func BasketExists(userId, basketId primitive.ObjectID) error {
	filter := bson.M{"_id": userId, "baskets": bson.M{"$elemMatch": bson.M{"_id": basketId}}}
	cursor := db.DB.Collection(constants.UsersCollection).FindOne(context.TODO(), filter)
	if cursor.Err() == mongo.ErrNoDocuments {
		return constants.ErrBasketNotFound
	}
	return cursor.Err()
}

func AddItemToBasket(userId, basketId primitive.ObjectID, item models.BasketItem) error {
	err := BasketExists(userId, basketId)
	if err != nil {
		return err
	}
	err = FindItemInBasket(item.ItemId, userId, basketId)
	if err != nil && err != constants.ErrItemNotFoundInBasket {
		return err
	}
	filter := bson.M{"_id": userId, "baskets": bson.M{"$elemMatch": bson.M{"_id": basketId}}}
	if err == constants.ErrItemNotFoundInBasket {
		update := bson.M{"$push": bson.M{"baskets.$.items": item}}
		_, err := db.DB.Collection(constants.UsersCollection).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		return nil
	}
	update := bson.M{"$set": bson.M{"baskets.$.items.$[i].qty": item.Qty}}
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"i.itemId": item.ItemId}},
	}
	opts := options.Update().SetArrayFilters(arrayFilters)
	_, err = db.DB.Collection(constants.UsersCollection).UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItemFromBasket(itemId, userId, basketId primitive.ObjectID) error {
	filter := bson.M{"_id": userId, "baskets": bson.M{"$elemMatch": bson.M{"_id": basketId}}}
	update := bson.M{"$pull": bson.M{"baskets.$.items": bson.M{"itemId": itemId}}}
	res, err := db.DB.Collection(constants.UsersCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return constants.ErrBasketNotFound
	}
	if res.ModifiedCount == 0 {
		return constants.ErrItemNotFoundInBasket
	}
	return nil
}
