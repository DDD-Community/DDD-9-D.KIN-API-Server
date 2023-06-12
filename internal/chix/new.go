package chix

import (
	"d.kin-app/internal/serverless"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() (r *chi.Mux) {
	r = chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	if serverless.IsLambdaRuntime() {
		r.Use(middleware.Recoverer)
	} else {
		r.Use(middleware.Logger)
	}
	r.Use(BodyNFCNormalize)
	return
}
