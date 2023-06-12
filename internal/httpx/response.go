package httpx

import (
	"encoding/json"
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
