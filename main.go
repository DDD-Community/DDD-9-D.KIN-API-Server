package main

import (
	"bytes"
	"context"
	"d.kin-app/internal/chix"
	"d.kin-app/internal/httpx"
	"d.kin-app/internal/serverless"
	"d.kin-app/internal/typex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"net/http"
)

type (
	apiGatewayHandler func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	lambdaHandler     func(context.Context, json.RawMessage) (json.RawMessage, error)
)

func wrapToWarmUp(fn apiGatewayHandler) lambdaHandler {
	var noData = []byte("{}")

	return func(ctx context.Context, payload json.RawMessage) (result json.RawMessage, err error) {
		if len(payload) == 2 && bytes.Equal(payload, noData) {
			fmt.Println("WARM UP")
			// to warm up
			result = noData
			return
		}

		var req events.APIGatewayV2HTTPRequest
		err = json.Unmarshal(payload, &req)
		if err != nil {
			return
		}

		res, err := fn(ctx, req)
		if err != nil {
			return
		}

		result, err = json.Marshal(res)
		return
	}
}

func main() {
	r := chix.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteContentType(w, httpx.ApplicationJSONCharsetUTF8)
		httpx.WriteStatus(w, http.StatusOK)
		httpx.WriteJSON(w, typex.JSONObject{
			"hello": "world",
		})
	})
	if serverless.IsLambdaRuntime() {
		lambda.Start(wrapToWarmUp(httpadapter.NewV2(r).ProxyWithContext))
	} else {
		http.ListenAndServe(":3000", r)
	}
}
