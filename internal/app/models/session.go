package models

import (
	"time"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Uid       primitive.ObjectID `bson:"uid"`
	Token     string             `bson:"token"`
	CreatedAt time.Time          `bson:"createdAt"`
	ExpiresAt time.Time          `bson:"expiresAt"`
}

func (*Session) Collection() string {
	return constants.SessionCollection
}
