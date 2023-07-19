package main

import (
	"context"
	"d.kin-app/internal/auth/fbauth/keystore"
	"d.kin-app/internal/auth/fbauth/token"
	"d.kin-app/internal/auth/util"
	"d.kin-app/internal/awsx"
	"d.kin-app/internal/awsx/lambdax"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	headerAuthorization = "authorization"
	iamEffectAllow      = "Allow"
	iamEffectDeny       = "Deny"
)

var (
	errUnauthorized = errors.New("Unauthorized")
)

var _ lambda.Handler = (*handler)(nil)

type handler struct {
	verifier *token.Verifier
}

func (h *handler) Invoke(ctx context.Context, payload []byte) (result []byte, err error) {
	if lambdax.IsEmptyPayload(payload) {
		fmt.Println("WARM UP") // TODO: structured logger 로변경 필요
		// to warm up
		result = lambdax.EmptyPayload
		return
	}

	var req events.APIGatewayV2CustomAuthorizerV2Request
	err = json.Unmarshal(payload, &req)
	if err != nil {
		return
	}

	rawJwt, err := util.BearerToken(req.Headers[headerAuthorization])
	if err != nil {
		err = errUnauthorized
		return
	}

	tk, err := h.verifier.Verify(rawJwt)
	if err != nil {
		err = errUnauthorized
		return
	}

	claims := tk.FirebaseClaims()

	// TODO: 유저 Get Or Create 구현
	var res events.APIGatewayV2CustomAuthorizerIAMPolicyResponse
	res.PrincipalID = claims.Subject
	res.Context = map[string]interface{}{
		// TODO: 유저 정보 전문 넘기기
		"user_id": claims.Subject,
	}
	res.PolicyDocument = allowPolicy()
	return json.Marshal(res)
}

func main() {
	lambda.Start(&handler{
		verifier: token.NewVerifier(
			keystore.NewKeyStoreByDynamoDB(
				awsx.DynamoDB.Value(),
			),
		),
	})
}

func allowPolicy() events.APIGatewayCustomAuthorizerPolicy {
	return generatePolicy(iamEffectAllow, "arn:aws:execute-api:ap-northeast-2:035366530565:*/*/*") // TODO: resource, 필요 하다면 명시적으로 변경
}

func generatePolicy(effect, resource string) events.APIGatewayCustomAuthorizerPolicy {
	if effect != "" && resource != "" {
		return events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	return events.APIGatewayCustomAuthorizerPolicy{}
}
