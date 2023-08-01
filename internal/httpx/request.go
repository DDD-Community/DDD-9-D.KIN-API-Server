package httpx

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RequestContentTypeIsString(r *http.Request) bool {
	return ContentTypeIsString(r.Header.Get("Content-Type"))
}

func ContentTypeIsString(contentType string) bool {
	return strings.HasPrefix(contentType, ApplicationJSON) ||
		strings.HasPrefix(contentType, ApplicationXML) ||
		strings.HasPrefix(contentType, TextXML) ||
		strings.HasPrefix(contentType, TextHTML) ||
		strings.HasPrefix(contentType, TextPlain)
}

func GetBodyJSON[T any](r *http.Request) (res T, err error) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), ApplicationJSON) {
		err = ErrContentTypeMustApplicationJSON
		return
	}

	err = json.NewDecoder(r.Body).Decode(&res)
	return
}

func DeviceId(r *http.Request) string {
	return r.Header.Get("Device-Id")
}
