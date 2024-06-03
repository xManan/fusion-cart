package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/xManan/fusion-cart/internal/app/utils"
	httpmux "github.com/xManan/fusion-cart/pkg/http-mux"
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
