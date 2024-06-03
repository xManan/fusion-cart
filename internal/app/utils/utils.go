package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func ErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	resp := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		ErrCode int    `json:"errCode"`
	}{
		false,
		message,
		code,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		ServerErrorResponse(w)
		return
	}
	fmt.Fprintf(w, "%s", string(jsonResp))
}

func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		true,
		message,
		data,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		ServerErrorResponse(w)
		return
	}
	fmt.Fprintf(w, "%s", string(jsonResp))
}

func BadRequestResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "%s", "Bad request")
}

func UnauthorizedResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "%s", "Unauthorized")
}

func ServerErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", "Something went wrong")
}

func MapRequestBodyToStruct(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		return err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return json.Unmarshal(body, v)
}

func IsUrl(url string) bool {
	if strings.Contains(url, " ") || !strings.Contains(url, "/") {
		return false
	}
	re := regexp.MustCompile(`\w+.\w+`)
	if re.MatchString(url) {
		return true
	}
	return false
}

func GenerateRandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
