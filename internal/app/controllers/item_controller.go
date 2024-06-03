package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/xManan/fusion-cart/internal/app/constants"
	"github.com/xManan/fusion-cart/internal/app/services"
	_ "github.com/xManan/fusion-cart/internal/app/services"
	"github.com/xManan/fusion-cart/internal/app/services/merchant"
	"github.com/xManan/fusion-cart/internal/app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleSearchItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println(err.Error())
		utils.ServerErrorResponse(w)
		return
	}
	if body.Query == "" {
		utils.ErrorResponse(w, "query is required", constants.CodeFieldRequired)
		return
	}
	if utils.IsUrl(strings.TrimSpace(body.Query)) {
		urlStr := body.Query
		if !strings.Contains(urlStr, "/") {
			utils.ErrorResponse(w, "Invalid URL. Try pasting a product link from supported stores", constants.CodeInvalidProductLink)
			return
		}
		urlStruct, err := url.Parse(urlStr)
		if err != nil {
			log.Println(err.Error())
			utils.ServerErrorResponse(w)
			return
		}

		link, err := merchant.NewProductLink(urlStruct)
		if err != nil {
			switch err {
			case merchant.ErrMerchantNotSupported:
				utils.ErrorResponse(w, "Store is not supported. Please try other stores", constants.CodeMerchantNotSupported)
				return
			case merchant.ErrItemRefNotFound:
				utils.ErrorResponse(w, "Invalid URL. Try pasting a product link from supported stores", constants.CodeInvalidProductLink)
				return
			default:
				log.Println(err.Error())
				utils.ServerErrorResponse(w)
				return
			}
		}
		item, err := services.FindItemByRef(link.ItemRef)
		if err == nil {
			utils.SuccessResponse(w, "", item)
			return
		} else if err != mongo.ErrNoDocuments {
			log.Println(err.Error())
			utils.ServerErrorResponse(w)
			return
		}

		item, err = merchant.FetchItem(&link)
		if err != nil {
			log.Println(err.Error())
			utils.ServerErrorResponse(w)
			return
		}
		err = services.SaveItem(&item)
		if err != nil {
			log.Println(err.Error())
			utils.ServerErrorResponse(w)
			return
		}
		utils.SuccessResponse(w, "", item)
	} else {
		utils.ErrorResponse(w, "Sorry, search functionality is currently limited to URLs only. Please paste a URL to search", constants.CodeSearchNotAvailable)
	}
}
