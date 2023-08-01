package dynamodbx

import (
	"context"
	"d.kin-app/internal/awsx"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetItem(ctx context.Context, tableName *string, key map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error) {
	return awsx.DynamoDB.Value().GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: tableName,
	})
}

func ExistItem(ctx context.Context, tableName *string, key map[string]types.AttributeValue) (bool, error) {
	res, err := GetItem(ctx, tableName, key)
	if err != nil {
		return false, err
	}

	return res.Item != nil, nil
}
