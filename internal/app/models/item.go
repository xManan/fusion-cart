package models

import (
	"time"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	Id        primitive.ObjectID `json:"id"        bson:"_id,omitempty"`
	Merchant  string             `json:"merchant"  bson:"merchant"`
	ItemRef   string             `json:"itemRef"   bson:"itemRef"`
	Name      string             `json:"name"      bson:"name"`
	Price     float64            `json:"price"     bson:"price"`
	OldPrice  float64            `json:"oldPrice"  bson:"oldPrice"`
	Url       string             `json:"url"       bson:"url"`
	ImgUrl    string             `json:"imgUrl"    bson:"imgUrl"`
	Rating    float64            `json:"rating"    bson:"rating"`
	MetaData  primitive.M        `json:"metaData"  bson:"metaData"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func(*Item) Collection() string {
	return constants.ItemsCollection
}
