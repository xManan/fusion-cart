package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"github.com/xManan/fusion-cart/internal/app/db"
	"github.com/xManan/fusion-cart/internal/app/models"
	"github.com/xManan/fusion-cart/internal/app/services"
	"github.com/xManan/fusion-cart/internal/app/types"
	"github.com/xManan/fusion-cart/internal/app/utils"
	"github.com/xManan/fusion-cart/internal/app/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	var user models.UnverifiedUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}

	if err := validators.ValidateRegistration(&user); err != nil {
		utils.ErrorResponse(w, err.Message, err.Code)
		return
	}

	emailRegistered, err := services.IsEmailRegistered(user.Email)
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	if emailRegistered {
		utils.ErrorResponse(w, "email is already registered", constants.CodeEmailAlreadyRegistered)
		return
	}

	if user.Mobile != "" {
		mobileInUse, err := services.IsMobileInUse(user.Mobile)
		if err != nil {
			log.Println(err.Error())
			utils.ServerErrorResponse(w)
			return
		}
		if mobileInUse {
			utils.ErrorResponse(w, "mobile is already in use", constants.CodeMobileAlreadyInUse)
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	user.Password = string(hash)

	token, err := services.CreateUnverifiedUser(&user)
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}

	_ = services.SendVerificationLink(user.Email, token)

	utils.SuccessResponse(w, "verification link sent", nil)
}

func HandleVerification(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.BadRequestResponse(w)
		return
	}
	var unverifiedUser models.UnverifiedUser
	err := db.DB.Collection(unverifiedUser.Collection()).FindOne(context.TODO(), bson.M{"token": token}).Decode(&unverifiedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.BadRequestResponse(w)
			return
		}
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}

	if time.Now().Compare(unverifiedUser.ExpiresAt) == 1 {
		utils.BadRequestResponse(w)
		return
	}

	err = services.MoveUnverifiedUserToUsers(&unverifiedUser)
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}

	utils.SuccessResponse(w, "user verified", nil)
	return
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body types.LoginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	if err := validators.ValidateLogin(&body); err != nil {
		utils.ErrorResponse(w, err.Message, err.Code)
		return
	}

	userId, err := services.VerifyCredentials(body.Email, body.Password)
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	if userId == "" {
		utils.ErrorResponse(w, "invalid email or password", constants.CodeInvalidEmailOrPassword)
		return
	}

	token, err := services.NewUserSession(userId)
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}

	utils.SuccessResponse(w, "logged in successfully", map[string]any{"token": token})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("bearerToken").(string)
	if !ok {
		log.Printf("Failed to convert token to string: %v\n", token)
		utils.ServerErrorResponse(w)
		return
	}

	res, err := db.DB.Collection(constants.SessionCollection).DeleteOne(context.TODO(), bson.M{ "token": token })
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	if res.DeletedCount == 0 {
		utils.UnauthorizedResponse(w)
		return
	}
	utils.SuccessResponse(w, "logged out successfully", nil)
}
