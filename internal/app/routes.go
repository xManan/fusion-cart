package app

import (
	"net/http"

	"github.com/xManan/fusion-cart/internal/app/controller"
	"github.com/xManan/fusion-cart/internal/app/middleware"
)

func RegisterRoutes(router *http.ServeMux) {

	router.HandleFunc("GET /{name}", middleware.Authenticate(controller.Index))

	router.HandleFunc("POST /search", controller.SearchItem)

}
