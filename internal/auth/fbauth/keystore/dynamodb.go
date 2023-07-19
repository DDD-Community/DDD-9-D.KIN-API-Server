package keystore

import (
	"context"
	"d.kin-app/internal/typex"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	dynamoDBTableName = "server"
	dynamoDBKey       = "fb_kid_cached"
)

func NewKeyStoreByDynamoDB(client *dynamodb.Client) KeyStore {
	return &basedDynamoDB{
		tableName: typex.P(dynamoDBTableName),
		key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: dynamoDBKey},
		},
		client: client,
	}
}

type basedDynamoDB struct {
	tableName *string
	key       map[string]types.AttributeValue
	client    *dynamodb.Client
}

func (db *basedDynamoDB) Get() (*PublicKeys, error) {
	resp, err := db.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		Key:       db.key,
		TableName: db.tableName,
	})
	if err != nil {
		return nil, err
	}
	if resp.Item == nil {
		return nil, ErrNotFound
	}

	var result struct {
		Kid       map[string]string `dynamodbav:"kid"`
		ExpiresAt int64             `dynamodbav:"expires_at"`
	}
	err = attributevalue.UnmarshalMap(resp.Item, &result)
	if err != nil {
		return nil, err
	}

	return &PublicKeys{
		KidMap:    result.Kid,
		ExpiresAt: result.ExpiresAt,
	}, nil
}

func (db *basedDynamoDB) Set(newKeys *PublicKeys) error {
	var data struct {
		Key       string            `dynamodbav:"key"`
		Kid       map[string]string `dynamodbav:"kid"`
		ExpiresAt int64             `dynamodbav:"expires_at"`
	}
	data.Key = dynamoDBKey
	data.Kid = newKeys.KidMap
	data.ExpiresAt = newKeys.ExpiresAt

	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		return err
	}

	_, err = db.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: db.tableName,
	})
	return err
}
