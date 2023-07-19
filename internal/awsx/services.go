package awsx

import (
	"d.kin-app/internal/config"
	"d.kin-app/internal/typex"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	DynamoDB = typex.ByLazy(func() *dynamodb.Client {
		return dynamodb.NewFromConfig(config.AWS.Value())
	})
)
