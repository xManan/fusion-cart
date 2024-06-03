package httpmux

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
	middlewares []func(next HttpHandler) HttpHandler
}

type HttpHandler = func(http.ResponseWriter, *http.Request)
type HttpMiddleware = func(next HttpHandler) HttpHandler

func NewRouter(mux *http.ServeMux) Router {
	return Router{ mux: mux }
}

func (router *Router) Use(middleware HttpMiddleware) {
	router.middlewares = append(router.middlewares, middleware)
}

func (router *Router) Route(method string, endpoint string, handler HttpHandler) {
	finalHandler := handler
	for _, middleware := range router.middlewares {
		finalHandler = middleware(finalHandler)
	}
	router.mux.HandleFunc(method + " " + endpoint, finalHandler)
}
