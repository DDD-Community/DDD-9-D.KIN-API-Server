package lambdax

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"net/http"
)

const (
	AuthContextKeyUserClaims = "user_claims"
)

func GetLambdaAuthorizer(r *http.Request) (_ *events.APIGatewayV2HTTPRequestContextAuthorizerDescription, ok bool) {
	ctx, ok := core.GetAPIGatewayV2ContextFromContext(r.Context())
	if !ok {
		return
	}

	return ctx.Authorizer, ctx.Authorizer != nil
}

func GetUserClaims(r *http.Request) (res map[string]any, ok bool) {
	a, ok := GetLambdaAuthorizer(r)
	if !ok {
		return
	}

	res, ok = a.Lambda[AuthContextKeyUserClaims].(map[string]any)
	return
}
