package main

import (
	"d.kin-app/internal/chix"
	"d.kin-app/internal/httpx"
	"d.kin-app/internal/typex"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"net/http"
)

func main() {
	r := chix.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteContentType(w, httpx.ApplicationJSONCharsetUTF8)
		httpx.WriteStatus(w, http.StatusOK)
		httpx.WriteJSON(w, typex.JSONObject{
			"hello": "world",
		})
	})
	lambda.Start(httpadapter.NewV2(r).ProxyWithContext)
}
