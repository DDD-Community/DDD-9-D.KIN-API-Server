package middleware

import (
	"context"
	"d.kin-app/internal/auth/fbauth/token"
	"d.kin-app/internal/awsx/lambdax"
	"net/http"
)

func UserClaimsExtractor(next http.Handler) http.Handler {
	var fn http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		uc, ok := lambdax.GetUserClaims(r)
		if ok {
			var tc token.Claims
			tc.Set(uc)
			ctx := context.WithValue(r.Context(), UserClaimsKey, &tc)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	}

	return fn
}

func GetUserClaims(r *http.Request) (res *token.Claims, ok bool) {
	res, ok = r.Context().Value(UserClaimsKey).(*token.Claims)
	return
}
