package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"github.com/xManan/fusion-cart/internal/app/db"
	"github.com/xManan/fusion-cart/internal/app/utils"
	httpmux "github.com/xManan/fusion-cart/pkg/http-mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ValidateJson(next httpmux.HttpHandler) httpmux.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "json") {
			body, err := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			if err != nil {
				utils.ServerErrorResponse(w)
				return
			}
			if !json.Valid(body) {
				utils.BadRequestResponse(w)
				return
			}
		}
		next(w, r)
	}
}

func Authenticate(next httpmux.HttpHandler) httpmux.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimLeft(authHeader, "Bearer ")
		if token == "" {
			utils.UnauthorizedResponse(w)
			return
		}
		var result struct {
			Id        primitive.ObjectID `bson:"_id"`
			Uid       primitive.ObjectID `bson:"uid"`
			ExpiresAt time.Time          `bson:"expiresAt"`
		}
		opts := options.FindOne().SetProjection(bson.M{"uid": 1, "expiresAt": 1})
		err := db.DB.Collection(constants.SessionCollection).FindOne(context.TODO(), bson.M{"token": token}, opts).Decode(&result)
		if err == mongo.ErrNoDocuments {
			utils.UnauthorizedResponse(w)
			return
		} else if err != nil {
			utils.ServerErrorResponse(w)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "bearerToken", token)
		ctx = context.WithValue(ctx, "userId", result.Uid)
		next(w, r.WithContext(ctx))
	}
}
