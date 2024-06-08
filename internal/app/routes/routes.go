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

	router.Route(http.MethodPost, "/api/v1/basket", middlewares.Authenticate(controllers.HandleAddBasket))
	router.Route(http.MethodGet, "/api/v1/basket", middlewares.Authenticate(controllers.HandleGetBaskets))
	router.Route(http.MethodGet, "/api/v1/basket/{id}", middlewares.Authenticate(controllers.HandleGetBasket))
	router.Route(http.MethodPut, "/api/v1/basket/{id}", middlewares.Authenticate(controllers.HandleEditBasket))
	router.Route(http.MethodDelete, "/api/v1/basket/{id}", middlewares.Authenticate(controllers.HandleDeleteBasket))
	router.Route(http.MethodPost, "/api/v1/basket/{id}/item", middlewares.Authenticate(controllers.HandleAddItem))
	router.Route(http.MethodDelete, "/api/v1/basket/{basketId}/item/{itemId}", middlewares.Authenticate(controllers.HandleDeleteItem))
}
