package user

import (
	"context"
	uioTypes "d.kin-app/cmd/lambda/userImgOptimizer/types"
	"d.kin-app/internal/awsx"
	"d.kin-app/internal/awsx/dynamodbx"
	"d.kin-app/internal/awsx/sqsx"
	"d.kin-app/internal/typex"
	"d.kin-app/internal/validation"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"strconv"
	"sync"
	"time"
)

const (
	DynamoDBKeyName   = "pk"
	DynamoDBTableName = "user"
)

var (
	//ptr_ddbKeyName   = typex.P(DynamoDBKeyName)
	ptr_ddbTableName         = typex.P(DynamoDBTableName)
	ptr_txCondExpNotExistsPK = typex.P(fmt.Sprintf("attribute_not_exists(%s)", DynamoDBKeyName))
)

type Gender string

const (
	GenderFemale Gender = "FEMALE"
	GenderMale   Gender = "MALE"
)

func GenderValues() []Gender {
	return []Gender{
		GenderFemale,
		GenderMale,
	}
}

type User struct {
	ctx           context.Context
	UserId        string            `dynamodbav:"pk"`
	Image         map[string]Image  `dynamodbav:"image,omitemptyelem"`
	ImageURL      *string           `dynamodbav:"image_url"`
	Nickname      string            `dynamodbav:"nickname"`
	YearOfBirth   int16             `dynamodbav:"year_of_birth"`
	Gender        Gender            `dynamodbav:"gender"`
	CreatedTime   int64             `dynamodbav:"created_time"`
	UpdatedTime   int64             `dynamodbav:"updated_time"`
	DeletedTime   *int64            `dynamodbav:"deleted_time,omitempty"`
	ArchiveTime   *int64            `dynamodbav:"archive_time,omitempty"`
	ActiveDevices map[string]Device `dynamodbav:"active_devices,omitemptyelem"`
}

func (u *User) SetContext(ctx context.Context) {
	u.ctx = ctx
}

func (u *User) Context() context.Context {
	if u.ctx == nil {
		u.SetContext(context.Background())
	}

	return u.ctx
}

var ptr_ddbUpdateExp_DoSignUp = typex.P(
	"SET nickname = :new_nickname, " +
		"year_of_birth = :new_year_of_birth, " +
		"gender = :new_gender, " +
		"updated_time = :new_updated_time",
)

func (u *User) DoSignUp(nickname string, yearOfBirth int16, gender Gender) (err error) {
	ctx := u.Context()
	err = IsUsableNickname(ctx, nickname)
	if err != nil {
		return
	}

	nowMilli := time.Now().UnixMilli()
	_, err = awsx.DynamoDB.Value().TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []ddbTypes.TransactWriteItem{
			{
				Update: &ddbTypes.Update{
					Key:                      u.ddbKey(),
					TableName:                ptr_ddbTableName,
					UpdateExpression:         ptr_ddbUpdateExp_DoSignUp,
					ExpressionAttributeNames: nil,
					ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
						":new_nickname":      &ddbTypes.AttributeValueMemberS{Value: nickname},
						":new_year_of_birth": &ddbTypes.AttributeValueMemberN{Value: strconv.Itoa(int(yearOfBirth))},
						":new_gender":        &ddbTypes.AttributeValueMemberS{Value: string(gender)},
						":new_updated_time":  &ddbTypes.AttributeValueMemberN{Value: strconv.FormatInt(nowMilli, 10)},
					},
					ReturnValuesOnConditionCheckFailure: ddbTypes.ReturnValuesOnConditionCheckFailureAllOld,
				},
			},
			{
				Put: &ddbTypes.Put{
					Item:                                makeDDBKeyNickname(nickname),
					TableName:                           ptr_ddbTableName,
					ConditionExpression:                 ptr_txCondExpNotExistsPK,
					ReturnValuesOnConditionCheckFailure: ddbTypes.ReturnValuesOnConditionCheckFailureAllOld,
				},
			},
		},
		ReturnConsumedCapacity:      ddbTypes.ReturnConsumedCapacityTotal,
		ReturnItemCollectionMetrics: ddbTypes.ReturnItemCollectionMetricsNone,
	})
	if err == nil {
		u.Nickname = nickname
		u.YearOfBirth = yearOfBirth
		u.Gender = gender
		u.UpdatedTime = nowMilli
	}
	return
}

var ptr_ddbUpdateExp_ProfileUpdate = typex.P(
	"SET image_url = :new_image_url, " +
		"nickname = :new_nickname, " +
		"updated_time = :new_updated_time",
)

func (u *User) ProfileUpdate(imageURL *string, nickname string) (err error) {
	isSameNickname := u.Nickname == nickname
	isSameImageURL := typex.PrimitiveDeepEq(u.ImageURL, imageURL)
	if isSameNickname && isSameImageURL {
		return
	}

	ctx := u.Context()

	nowMilli := time.Now().UnixMilli()
	txItems := make([]ddbTypes.TransactWriteItem, 0, 3)
	if !isSameNickname {
		err = IsUsableNickname(ctx, nickname)
		if err != nil {
			return
		}

		txItems = append(txItems, ddbTypes.TransactWriteItem{
			Put: &ddbTypes.Put{
				Item:                                makeDDBKeyNickname(nickname),
				TableName:                           ptr_ddbTableName,
				ConditionExpression:                 ptr_txCondExpNotExistsPK,
				ReturnValuesOnConditionCheckFailure: ddbTypes.ReturnValuesOnConditionCheckFailureAllOld,
			},
		})

		oldNickname := u.Nickname
		defer func() {
			if err == nil {
				awsx.DynamoDB.Value().DeleteItem(ctx, &dynamodb.DeleteItemInput{
					Key:       makeDDBKeyNickname(oldNickname),
					TableName: ptr_ddbTableName,
				})
			}
		}()
	}

	if !isSameImageURL && imageURL != nil {

		defer func() {
			if err == nil {
				sqsx.SendData(
					ctx,
					"https://sqs.ap-northeast-2.amazonaws.com/035366530565/image-optimize",
					uioTypes.UserImageOptimize{
						UserId:  u.UserId,
						ImageId: imageId(*imageURL),
					},
				)
			}
		}()
	}

	txItems = append(txItems, ddbTypes.TransactWriteItem{
		Update: &ddbTypes.Update{
			Key:              u.ddbKey(),
			TableName:        ptr_ddbTableName,
			UpdateExpression: ptr_ddbUpdateExp_ProfileUpdate,
			ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
				":new_image_url":    dynamodbx.StringOrNull(imageURL),
				":new_nickname":     dynamodbx.String(nickname),
				":new_updated_time": dynamodbx.Signed(nowMilli),
			},
			ReturnValuesOnConditionCheckFailure: ddbTypes.ReturnValuesOnConditionCheckFailureAllOld,
		},
	})

	_, err = awsx.DynamoDB.Value().TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems:               txItems,
		ReturnConsumedCapacity:      ddbTypes.ReturnConsumedCapacityTotal,
		ReturnItemCollectionMetrics: ddbTypes.ReturnItemCollectionMetricsNone,
	})
	if err != nil {
		return
	}

	u.ImageURL = imageURL
	u.Nickname = nickname
	u.UpdatedTime = nowMilli
	return
}

var ptr_ddbUpdateExp_MakeProfileImageUploadURL = typex.P("SET image = :new_image")

func (u *User) MakeProfileImageUploadURL(file ImageFile) (res UploadLink, err error) {
	img := makeImage(file)
	ctx := u.Context()
	resp, err := awsx.S3Presign.Value().PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        &img.S3Bucket,
		Key:           &img.S3ObjectKey,
		ContentLength: img.File.Size,
		ContentType:   &img.File.MimeType,
	})
	if err != nil {
		return
	}

	newImage := u.cloneImage()
	newImage[img.ImageId] = img

	_, err = awsx.DynamoDB.Value().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:       u.ddbKey(),
		TableName: ptr_ddbTableName,
		ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
			":new_image": typex.Must(attributevalue.Marshal(newImage)),
		},
		ReturnValuesOnConditionCheckFailure: ddbTypes.ReturnValuesOnConditionCheckFailureAllOld,
		UpdateExpression:                    ptr_ddbUpdateExp_MakeProfileImageUploadURL,
	})
	if err != nil {
		return
	}

	u.Image = newImage

	res.URL = resp.URL
	res.Method = resp.Method
	return
}

var ptr_ddbUpdateExp_OptimizeImage = typex.P(
	"SET image = :new_image, " +
		"image_url = :new_image_url",
)

func (u *User) OptimizeImage(imageId string) {
	oldImage, ok := u.Image[imageId]
	if !ok {
		return
	}

	ctx := u.Context()
	optimizedImage, err := oldImage.makeWebP(ctx)
	if err != nil {
		return
	}

	_, err = awsx.DynamoDB.Value().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:       u.ddbKey(),
		TableName: ptr_ddbTableName,
		ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
			":new_image": typex.Must(attributevalue.Marshal(map[string]Image{
				optimizedImage.ImageId: optimizedImage,
			})),
			"new_image_url": dynamodbx.String(optimizedImage.imageURL()),
		},
		UpdateExpression: ptr_ddbUpdateExp_OptimizeImage,
	})

	var wg sync.WaitGroup
	deleteImages := make(map[string][]s3Types.ObjectIdentifier)
	for _, v := range u.Image {
		deleteImages[v.S3Bucket] = append(deleteImages[v.S3Bucket], s3Types.ObjectIdentifier{
			Key: &v.S3ObjectKey,
		})
	}

	for k, v := range deleteImages {
		bucket := k
		objects := v
		go func() {
			defer wg.Done()
			awsx.S3.Value().DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: &bucket,
				Delete: &s3Types.Delete{
					Objects: objects,
					Quiet:   true,
				},
			})
		}()
	}

	wg.Wait()
}

func (u *User) cloneImage() (res map[string]Image) {
	if u.Image == nil {
		return
	}

	res = make(map[string]Image)
	for k, v := range u.Image {
		res[k] = v
	}
	return
}

func (u *User) ddbKey() map[string]ddbTypes.AttributeValue {
	return dynamodbx.SingleStringField(DynamoDBKeyName, u.UserId)
}

func (u *User) ddbKeyNickname() map[string]ddbTypes.AttributeValue {
	return makeDDBKeyNickname(u.Nickname)
}

func IsUsableNickname(ctx context.Context, nickname string) (err error) {
	if !validation.IsValidNickname(nickname) {
		return ErrInvalidNickname
	}

	nicknameExists, err := dynamodbx.ExistItem(ctx, ptr_ddbTableName, makeDDBKeyNickname(nickname))
	if err != nil {
		return
	} else if nicknameExists {
		return ErrNicknameAlreadyExists
	}

	return
}

func GetUser(ctx context.Context, key string) (_ *User, err error) {
	res, err := getItem(ctx, key)
	if err != nil {
		return
	}
	if res.Item == nil {
		err = ErrUserNotFound
		return
	}

	var u User
	err = attributevalue.UnmarshalMap(res.Item, &u)
	if err != nil {
		return
	}
	u.SetContext(ctx)

	return &u, nil
}

func CreateUser(ctx context.Context, u *User) (err error) {
	txItems := make([]ddbTypes.TransactWriteItem, 0, 2)
	if len(u.Nickname) > 0 {
		uerr := IsUsableNickname(ctx, u.Nickname)
		if uerr != nil {
			return uerr
		}

		txItems = append(txItems, ddbTypes.TransactWriteItem{
			Put: &ddbTypes.Put{
				Item:                                u.ddbKeyNickname(),
				TableName:                           ptr_ddbTableName,
				ConditionExpression:                 ptr_txCondExpNotExistsPK,
				ReturnValuesOnConditionCheckFailure: ddbTypes.ReturnValuesOnConditionCheckFailureAllOld,
			},
		})
	}

	item, err := attributevalue.MarshalMap(u)
	if err != nil {
		return
	}

	txItems = append(txItems, ddbTypes.TransactWriteItem{
		Put: &ddbTypes.Put{
			Item:                                item,
			TableName:                           ptr_ddbTableName,
			ConditionExpression:                 ptr_txCondExpNotExistsPK,
			ReturnValuesOnConditionCheckFailure: ddbTypes.ReturnValuesOnConditionCheckFailureAllOld,
		},
	})

	_, err = awsx.DynamoDB.Value().TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems:               txItems,
		ReturnConsumedCapacity:      ddbTypes.ReturnConsumedCapacityTotal,
		ReturnItemCollectionMetrics: ddbTypes.ReturnItemCollectionMetricsNone,
	})
	if err == nil {
		u.SetContext(ctx)
	}
	return
}

func makeDDBKeyNickname(nickname string) map[string]ddbTypes.AttributeValue {
	return dynamodbx.SingleStringField(
		DynamoDBKeyName,
		formatNicknameKey(nickname),
	)
}

func formatNicknameKey(nickname string) string {
	return fmt.Sprintf("nickname#%s", nickname)
}

func getItem(ctx context.Context, key string) (*dynamodb.GetItemOutput, error) {
	return dynamodbx.GetItem(ctx, ptr_ddbTableName, dynamodbx.SingleStringField(
		DynamoDBKeyName,
		key,
	))
}
