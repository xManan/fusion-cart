package controller

import (
	"net/http"

	"github.com/xManan/fusion-cart/internal/app/validator"
)

func SearchItem(w http.ResponseWriter, r *http.Request) {
	err := validator.SearchItemRequest(r)
	if err != nil {
		
	}
}
