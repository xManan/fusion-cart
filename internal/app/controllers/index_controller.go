package controllers

import (
	"fmt"
	"net/http"
)


func HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Working fine")
}
