package lambdax

import (
	"bytes"
	"os"
)

const (
	envLambdaServerPort = "_LAMBDA_SERVER_PORT"
	envLambdaRuntimeAPI = "AWS_LAMBDA_RUNTIME_API"
)

func IsLambdaRuntime() bool {
	return isLambdaRuntime()
}

func isLambdaRuntime() bool {
	_, ok1 := os.LookupEnv(envLambdaServerPort)
	_, ok2 := os.LookupEnv(envLambdaRuntimeAPI)
	return ok1 && ok2
}

var (
	EmptyPayload = []byte("{}")
)

func IsEmptyPayload(payload []byte) bool {
	return len(payload) == 2 && bytes.Equal(payload, EmptyPayload)
}
