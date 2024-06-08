package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"github.com/xManan/fusion-cart/internal/app/models"
	"github.com/xManan/fusion-cart/internal/app/services"
	"github.com/xManan/fusion-cart/internal/app/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleAddBasket(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}

	userId, ok := r.Context().Value("userId").(primitive.ObjectID)
	if !ok {
		log.Printf("Failed to convert userId to primitive.ObjectID: %v\n", userId)
		utils.ServerErrorResponse(w)
		return
	}

	err := services.AddBasketToUserBaskets(userId, body.Name)
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	utils.SuccessResponse(w, "basket added successfully", nil)
}

func HandleGetBaskets(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(primitive.ObjectID)
	if !ok {
		log.Printf("Failed to convert userId to primitive.ObjectID: %v\n", userId)
		utils.ServerErrorResponse(w)
		return
	}
	baskets, err := services.GetBaskets(userId)
	if err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	utils.SuccessResponse(w, "", map[string]interface{}{
		"baskets": baskets,
	})
}

func HandleGetBasket(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(primitive.ObjectID)
	if !ok {
		log.Printf("Failed to convert userId to primitive.ObjectID: %v\n", userId)
		utils.ServerErrorResponse(w)
		return
	}
	basketId := r.PathValue("id")
	basketOid, err := primitive.ObjectIDFromHex(basketId)
	if err != nil {
		utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
		return
	}
	basket, err := services.GetBasket(userId, basketOid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
			return
		}
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	utils.SuccessResponse(w, "", map[string]interface{}{
		"basket": basket,
	})
}

func HandleEditBasket(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		if err == io.EOF {
			utils.BadRequestResponse(w)
			return
		}
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	if body.Name == "" {
		utils.ErrorResponse(w, "basket name is required", constants.CodeFieldRequired)
		return
	}

	userId, ok := r.Context().Value("userId").(primitive.ObjectID)
	if !ok {
		log.Printf("Failed to convert userId to primitive.ObjectID: %v\n", userId)
		utils.ServerErrorResponse(w)
		return
	}

	basketId := r.PathValue("id")
	basketOid, err := primitive.ObjectIDFromHex(basketId)
	if err != nil {
		utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
		return
	}

	err = services.UpdateBasket(userId, basketOid, body.Name)
	if err != nil {
		if err == constants.ErrBasketNotFound {
			utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
			return
		}
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	utils.SuccessResponse(w, "basket updated successfully", nil)
}

func HandleDeleteBasket(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(primitive.ObjectID)
	if !ok {
		log.Printf("Failed to convert userId to primitive.ObjectID: %v\n", userId)
		utils.ServerErrorResponse(w)
		return
	}

	basketId := r.PathValue("id")
	basketOid, err := primitive.ObjectIDFromHex(basketId)
	if err != nil {
		utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
		return
	}
	err = services.DeleteBasket(userId, basketOid)
	if err != nil {
		if err == constants.ErrBasketNotFound {
			utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
			return
		}
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	utils.SuccessResponse(w, "basket deleted successfully", nil)
}

func HandleAddItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ItemId string `json:"itemId"`
		Qty    uint64 `json:"qty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	if body.Qty <= 0 {
		utils.ErrorResponse(w, "qty must be greater than zero", constants.CodeQtyMustBeGreaterThanZero)
		return
	}

	itemOid, err := primitive.ObjectIDFromHex(body.ItemId)
	if err != nil {
		utils.ErrorResponse(w, "item not found", constants.CodeItemNotFound)
		return
	}

	item := models.BasketItem{ItemId: itemOid, Qty: body.Qty}
	userId := r.Context().Value("userId").(primitive.ObjectID)

	basketId := r.PathValue("id")
	basketOid, err := primitive.ObjectIDFromHex(basketId)
	if err != nil {
		utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
		return
	}

	err = services.AddItemToBasket(userId, basketOid, item)
	if err != nil {
		switch err {
		case constants.ErrBasketNotFound:
			utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
			return
		case constants.ErrItemNotFoundInBasket:
			utils.ErrorResponse(w, "item not found in basket", constants.CodeItemNotFoundInBasket)
			return
		}
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	utils.SuccessResponse(w, "item added successfully", nil)
}

func HandleDeleteItem(w http.ResponseWriter, r *http.Request) {
	basketId := r.PathValue("basketId")
	basketOid, err := primitive.ObjectIDFromHex(basketId)
	if err != nil {
		utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
	}
	itemId := r.PathValue("itemId")
	itemOid, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
	}
	userId := r.Context().Value("userId").(primitive.ObjectID)
	err = services.DeleteItemFromBasket(itemOid, userId, basketOid)
	if err != nil {
		switch err {
		case constants.ErrBasketNotFound:
			utils.ErrorResponse(w, "basket not found", constants.CodeBasketNotFound)
			return
		case constants.ErrItemNotFoundInBasket:
			utils.ErrorResponse(w, "item not found in basket", constants.CodeItemNotFoundInBasket)
			return
		}
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	utils.SuccessResponse(w, "item deleted successfully", nil)
}
