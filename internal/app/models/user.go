package models

import (
	"time"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Email     string             `bson:"email"`
	Mobile    string             `bson:"mobile"`
	Password  string             `bson:"password"`
	Baskets   []Basket           `bson:"baskets"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func (*User) Collection() string {
	return constants.UsersCollection
}
