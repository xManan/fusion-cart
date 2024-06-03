package models

import (
	"time"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnverifiedUser struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Email     string             `bson:"email"`
	Mobile    string             `bson:"mobile"`
	Password  string             `bson:"password"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expiresAt"`
}

func (*UnverifiedUser) Collection() string {
	return constants.UnverifiedUsersCollection
}
