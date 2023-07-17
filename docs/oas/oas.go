package oas

import (
	"d.kin-app/internal/httpx"
	_ "embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"os"
)

//go:embed gen/openapi/openapi.yaml
var docs []byte

func Handler(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteContentType(w, "text/html; charset=UTF-8")
	//httpx.WriteContentType(w, "text/yaml; charset=UTF-8")
	Write(w)
}

func Write(w io.Writer) (int, error) {
	return w.Write(docs)
}

//go:embed oas-ui.html
var ui []byte

func UIHandler(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteContentType(w, "text/html; charset=UTF-8")
	//httpx.WriteContentType(w, "text/yaml; charset=UTF-8")
	UIWrite(w)
}

func UIWrite(w io.Writer) (int, error) {
	return w.Write(ui)
}

func Route(r chi.Router) {
	r.Route("/oas/docs", func(r chi.Router) {
		r.Get("/", UIHandler)
		r.Get("/openapi.yaml", Handler)
	})
}

const (
	EnvOASUsername = "OAS_USERNAME"
	EnvOASPassword = "OAS_PASSWORD"
)

func RouteWithBasicAuth(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth("API Docs", map[string]string{
			os.Getenv(EnvOASUsername): os.Getenv(EnvOASPassword),
		}))
		Route(r)
	})
}
