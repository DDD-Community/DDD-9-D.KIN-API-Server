package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func WriteContentType(w http.ResponseWriter, mimeType string) {
	w.Header().Set("Content-Type", mimeType)
}

func WriteStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func WriteJSON(w http.ResponseWriter, data any) error {
	return json.NewEncoder(w).Encode(data)
}

func SetBodyJSON(w http.ResponseWriter, status int, data any) error {
	WriteContentType(w, ApplicationJSONCharsetUTF8)
	WriteStatus(w, status)
	return WriteJSON(w, data)
}

func Error(w http.ResponseWriter, err error, status int) {
	var data struct {
		Message string `json:"message"`
	}
	data.Message = fmt.Sprintln(err)
	_ = SetBodyJSON(w, status, data)
}

func ErrorBadRequest(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusBadRequest)
}

var errInternalServerError = errors.New("internal server application error")

func ErrorInternalServerError(w http.ResponseWriter) {
	Error(w, errInternalServerError, http.StatusInternalServerError)
}

var errUnauthorized = errors.New("Unauthorized")

func ErrorUnauthorized(w http.ResponseWriter) {
	Error(w, errUnauthorized, http.StatusUnauthorized)
}

type APIError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (err APIError) Error() string {
	return fmt.Sprintf("http status=%d, code=%s, message=%s", err.HTTPStatus, err.Code, err.Message)
}

func (err APIError) ResponseWrite(w http.ResponseWriter) error {
	return SetBodyJSON(w, err.HTTPStatus, err)
}
