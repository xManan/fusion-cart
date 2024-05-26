package middleware

import "net/http"

func Authenticate(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if true {
			return
		}
		next(w, r)
	}
}
