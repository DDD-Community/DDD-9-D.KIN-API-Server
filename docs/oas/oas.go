package oas

import (
	"d.kin-app/internal/httpx"
	"net/http"
)

// TODO: embed file
var docs []byte

func HandleOAS(w http.ResponseWriter, r *http.Request) {
	httpx.WriteContentType(w, "text/yaml")
	_, err := w.Write(docs)
	if err != nil {
		// TODO: Logging, report
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
