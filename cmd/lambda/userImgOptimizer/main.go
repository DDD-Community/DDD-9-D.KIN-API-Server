package main

import (
	"context"
	"d.kin-app/cmd/lambda/userImgOptimizer/types"
	"d.kin-app/models/user"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handle(ctx context.Context, input events.SQSEvent) string {
	var data types.UserImageOptimize
	for i := range input.Records {
		err := json.Unmarshal([]byte(input.Records[i].Body), &data)
		if err != nil {
			continue
		}

		u, _ := user.GetUser(ctx, data.UserId)
		if u == nil {
			continue
		}

		u.OptimizeImage(data.ImageId)
	}
	return "ok"
}

func main() {
	lambda.Start(handle)
}
