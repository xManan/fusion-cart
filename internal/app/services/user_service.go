package services

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"github.com/xManan/fusion-cart/internal/app/db"
	"github.com/xManan/fusion-cart/internal/app/models"
	"github.com/xManan/fusion-cart/internal/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/crypto/bcrypt"
)

func IsEmailRegistered(email string) (bool, error) {
	res := db.DB.Collection(constants.UsersCollection).FindOne(context.TODO(), bson.M{"email": email})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsMobileInUse(mobile string) (bool, error) {
	res := db.DB.Collection(constants.UsersCollection).FindOne(context.TODO(), bson.M{"mobile": mobile})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateUnverifiedUser(u *models.UnverifiedUser) (string, error) {
	token, err := utils.GenerateRandomToken(64)
	if err != nil {
		return "", err
	}
	u.Token = token
	tokenExpiry := os.Getenv("TOKEN_EXPIRY")
	tokenExpiryInt, err := strconv.Atoi(tokenExpiry)
	if err != nil {
		tokenExpiryInt = 15
	}
	u.ExpiresAt = time.Now().Add(time.Duration(tokenExpiryInt) * time.Minute)
	_, err = db.DB.Collection(u.Collection()).InsertOne(context.TODO(), u)
	if err != nil {
		return "", err
	}
	return token, nil
}

// TODO
func SendVerificationLink(email string, token string) error {
	return nil
}

func MoveUnverifiedUserToUsers(u *models.UnverifiedUser) error {
	user := models.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Mobile:    u.Mobile,
		Password:  u.Password,
		Baskets:   []models.Basket{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	callback := func(sesctx mongo.SessionContext) (interface{}, error) {
		_, err := db.DB.Collection(constants.UsersCollection).InsertOne(sesctx, user)
		if err != nil {
			return nil, err
		}
		_, err = db.DB.Collection(constants.UnverifiedUsersCollection).DeleteOne(sesctx, bson.M{"_id": u.Id})
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := db.DB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), callback, txnOptions)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func VerifyCredentials(email string, password string) (string, error) {
	var result struct {
		Id       primitive.ObjectID `bson:"_id"`
		Email    string             `bson:"email"`
		Password string             `bson:"password"`
	}
	options := options.FindOne().SetProjection(bson.M{"email": 1, "password": 1})
	err := db.DB.Collection(constants.UsersCollection).FindOne(context.TODO(), bson.M{"email": email}, options).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)) != nil {
		return "", nil
	}
	return result.Id.Hex(), nil
}

func NewUserSession(userId string) (string, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateRandomToken(64)
	if err != nil {
		return "", err
	}
	sessionExpiry := os.Getenv("SESSION_EXPIRY")
	sessionExpiryInt, err := strconv.Atoi(sessionExpiry)
	if err != nil {
		sessionExpiryInt = 12 * 60
	}
	expiresAt := time.Now().Add(time.Duration(sessionExpiryInt) * time.Minute)
	session := models.Session{
		Uid:       oid,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
	_, err = db.DB.Collection(session.Collection()).InsertOne(context.TODO(), session)
	if err != nil {
		return "", err
	}
	return token, nil
}
