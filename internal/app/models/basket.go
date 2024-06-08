package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BasketItem struct {
	ItemId primitive.ObjectID `json:"itemId" bson:"itemId"`
	Qty    uint64             `json:"qty"    bson:"qty"`
}

type Basket struct {
	Id    primitive.ObjectID `json:"id"    bson:"_id,omitempty"`
	Name  string             `json:"name"  bson:"name"`
	Items []BasketItem       `json:"items" bson:"items"`
}
