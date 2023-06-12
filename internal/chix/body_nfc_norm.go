package chix

import (
	"d.kin-app/internal/httpx"
	"golang.org/x/text/unicode/norm"
	"io"
	"net/http"
)

func BodyNFCNormalize(next http.Handler) http.Handler {
	var fn http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		if httpx.RequestContentTypeIsString(r) && r.Body != nil {
			r.Body = io.NopCloser(norm.NFC.Reader(r.Body))
		}
		next.ServeHTTP(w, r)
	}

	return fn
}
