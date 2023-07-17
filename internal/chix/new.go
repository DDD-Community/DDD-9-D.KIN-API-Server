package chix

import (
	"d.kin-app/docs/oas"
	"d.kin-app/internal/serverless"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() (r *chi.Mux) {
	r = chi.NewRouter()
	if serverless.IsLambdaRuntime() {
		r.Use(middleware.Recoverer)
	}

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		BodyNFCNormalize,
	)

	oas.RouteWithBasicAuth(r)
	return
}
