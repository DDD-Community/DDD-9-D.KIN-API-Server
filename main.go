package main

import (
	"context"
	"d.kin-app/internal/awsx"
	"d.kin-app/internal/awsx/lambdax"
	"d.kin-app/internal/chix"
	"d.kin-app/internal/httpx"
	"d.kin-app/internal/typex"
	"d.kin-app/routers/user"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"net/http"
)

var _ lambda.Handler = (*handler)(nil)

type handler struct {
	adapter *httpadapter.HandlerAdapterV2
}

func (h *handler) Invoke(ctx context.Context, payload []byte) (result []byte, err error) {
	if lambdax.IsEmptyPayload(payload) {
		fmt.Println("WARM UP") // TODO: structured logger 로변경 필요
		// to warm up
		result = lambdax.EmptyPayload
		return
	}

	var req events.APIGatewayV2HTTPRequest
	err = json.Unmarshal(payload, &req)
	if err != nil {
		return
	}

	res, err := h.adapter.ProxyWithContext(ctx, req)
	if err != nil {
		return
	}

	return json.Marshal(res)
}

// TODO: to be deleted
func createHelloWorldHandler(key, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteContentType(w, httpx.ApplicationJSONCharsetUTF8)
		httpx.WriteStatus(w, http.StatusOK)
		httpx.WriteJSON(w, typex.JSONObject{
			key: value,
		})
	}
}

func main() {
	resp, err := awsx.S3.Value().DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
		Bucket: typex.P("dkin-attachment"),
		Delete: &types.Delete{
			Objects: []types.ObjectIdentifier{
				{
					Key: typex.P("local/2e3eb1c0-8efb-41b5-9879-9133e21a2805"),
				},
				{
					Key: typex.P("local/2e3eb1c0-8efb-41b5-9879-9133e21a2805"),
				},
				{
					Key: typex.P("asdfqwe"),
				},
			},
			Quiet: true,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Errors)
	fmt.Println(resp.Errors)
	r := chix.NewRouter()
	r.Get("/", createHelloWorldHandler("hello", "world"))
	r.Get("/need-auth", createHelloWorldHandler("hi", "user"))
	r.Mount("/", user.NewHTTPHandler())
	if lambdax.IsLambdaRuntime() {
		lambda.Start(&handler{
			adapter: httpadapter.NewV2(r),
		})
	} else {
		http.ListenAndServe(":3000", r)
	}
}
