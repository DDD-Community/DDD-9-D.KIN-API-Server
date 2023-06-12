package httpx

import (
	"net/http"
	"strings"
)

func RequestContentTypeIsString(r *http.Request) bool {
	return ContentTypeIsString(r.Header.Get("Content-Type"))
}

func ContentTypeIsString(contentType string) bool {
	return strings.HasPrefix(contentType, "application/json") ||
		strings.HasPrefix(contentType, "application/xml") ||
		strings.HasPrefix(contentType, "text/xml") ||
		strings.HasPrefix(contentType, "text/html") ||
		strings.HasPrefix(contentType, "text/plain")
}
