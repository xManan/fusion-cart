package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BasketItem struct {
	ItemId primitive.ObjectID `json:"itemId" bson:"itemId"`
	Qty    uint64             `json:"qty" bson:"qty"`
}

type Basket struct {
	Name  string       `json:"name" bson:"name"`
	Label string       `json:"label" bson:"label"`
	Items []BasketItem `json:"items" bson:"items"`
}
