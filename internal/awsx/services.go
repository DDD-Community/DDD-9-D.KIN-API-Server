package awsx

import (
	"d.kin-app/internal/config"
	"d.kin-app/internal/typex"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var (
	DynamoDB = typex.ByLazy(func() *dynamodb.Client {
		return dynamodb.NewFromConfig(config.AWS.Value())
	})

	S3 = typex.ByLazy(func() *s3.Client {
		return s3.NewFromConfig(config.AWS.Value())
	})

	S3Presign = typex.ByLazy(func() *s3.PresignClient {
		return s3.NewPresignClient(S3.Value())
	})

	SQS = typex.ByLazy(func() *sqs.Client {
		return sqs.NewFromConfig(config.AWS.Value())
	})
)
