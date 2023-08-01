package sqsx

import (
	"context"
	"d.kin-app/internal/awsx"
	"d.kin-app/internal/typex"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SendData(ctx context.Context, queueURL string, data any) (err error) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}

	_, err = awsx.SQS.Value().SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: typex.P(string(msg)),
		QueueUrl:    &queueURL,
	})
	return
}
