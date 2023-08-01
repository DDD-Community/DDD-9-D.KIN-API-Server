package dynamodbx

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/exp/constraints"
	"strconv"
)

func String(str string) types.AttributeValue {
	return &types.AttributeValueMemberS{Value: str}
}

func StringOrNull(str *string) types.AttributeValue {
	if str != nil {
		return String(*str)
	}

	return &types.AttributeValueMemberNULL{Value: true}
}

func Signed[T constraints.Signed](v T) types.AttributeValue {
	return &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(v), 10)}
}

func Unsigned[T constraints.Unsigned](v T) types.AttributeValue {
	return &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(v), 10)}
}

func SingleStringField(keyName string, keyValue string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		keyName: &types.AttributeValueMemberS{Value: keyValue},
	}
}
