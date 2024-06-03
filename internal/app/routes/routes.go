package routes

import (
	"net/http"

	"github.com/xManan/fusion-cart/internal/app/controllers"
	"github.com/xManan/fusion-cart/internal/app/middlewares"
	_ "github.com/xManan/fusion-cart/internal/app/middlewares"
	httpmux "github.com/xManan/fusion-cart/pkg/http-mux"
)

func RegisterRoutes(router *httpmux.Router) {
	router.Use(middlewares.ValidateJson)

	router.Route(http.MethodGet, "/healthcheck", controllers.HandleHealthcheck)

	router.Route(http.MethodPost, "/api/v1/search", controllers.HandleSearchItem)

	router.Route(http.MethodPost, "/api/v1/register", controllers.HandleRegistration)
	router.Route(http.MethodGet, "/api/v1/verify", controllers.HandleVerification)
	router.Route(http.MethodPost, "/api/v1/login", controllers.HandleLogin)
	router.Route(http.MethodPut, "/api/v1/logout", middlewares.Authenticate(controllers.HandleLogout))
}
