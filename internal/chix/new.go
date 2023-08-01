package chix

import (
	"d.kin-app/docs/oas"
	"d.kin-app/internal/awsx/lambdax"
	serviceMiddleware "d.kin-app/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() (r *chi.Mux) {
	r = chi.NewRouter()
	if lambdax.IsLambdaRuntime() {
		r.Use(middleware.Recoverer)
	}

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		serviceMiddleware.BodyNFCNormalize,
		serviceMiddleware.UserClaimsExtractor,
	)

	oas.RouteWithBasicAuth(r)
	return
}
