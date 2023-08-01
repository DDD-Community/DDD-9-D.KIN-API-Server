package user

import (
	"context"
	"d.kin-app/internal/awsx"
	"d.kin-app/internal/awsx/lambdax"
	"d.kin-app/internal/sha3x"
	"d.kin-app/internal/typex"
	"d.kin-app/modules/ffmpeg"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"path"
	"time"
)

const (
	imageS3Bucket = "dkin-attachment"
)

var (
	imageS3ObjectPrefix = typex.ByLazy(func() string {
		if lambdax.IsLambdaRuntime() {
			return "live"
		}
		return "local"
	})
)

type Image struct {
	ImageId           string    `dynamodbav:"-"`
	S3Bucket          string    `dynamodbav:"s3_bucket"`
	S3ObjectKey       string    `dynamodbav:"s3_object_key"`
	S3UploadURL       *string   `dynamodbav:"s3_upload_url,omitempty"`
	S3UploadMethod    *string   `dynamodbav:"s3_upload_method,omitempty"`
	S3UploadExpiresAt *int64    `dynamodbav:"s3_upload_expires_at,omitempty"`
	File              ImageFile `dynamodbav:"file"`
}

func (i *Image) genImageId() {
	i.ImageId = imageId(i.ImageURL())
}

func (i *Image) ImageURL() string {
	return fmt.Sprintf("https://%s.s3.ap-northeast-2.amazonaws.com/%s", i.S3Bucket, i.S3ObjectKey)
}

func (i *Image) makeWebP(ctx context.Context) (_ Image, err error) {
	resp, err := awsx.S3.Value().GetObject(ctx, &s3.GetObjectInput{
		Bucket: &i.S3Bucket,
		Key:    &i.S3ObjectKey,
	})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := ffmpeg.EncodeWebP(ctx, resp.Body)
	if err != nil {
		return
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		return
	}

	optimizedImage := Image{
		S3Bucket:    i.S3Bucket,
		S3ObjectKey: fmt.Sprintf("%s.webp", i.S3ObjectKey),
		File: ImageFile{
			Size:     fileStat.Size(),
			MimeType: "image/webp",
		},
	}
	optimizedImage.genImageId()

	_, err = awsx.S3.Value().PutObject(ctx, &s3.PutObjectInput{
		Bucket: &optimizedImage.S3Bucket,
		Key:    &optimizedImage.S3ObjectKey,
		Body:   file,
	})
	if err != nil {
		return
	}

	return optimizedImage, nil
}

func (i *Image) makeUploadURL(ctx context.Context) (err error) {
	exp := time.Now().Add(time.Hour)
	resp, err := awsx.S3Presign.Value().PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        &i.S3Bucket,
		Key:           &i.S3ObjectKey,
		ContentLength: i.File.Size,
		ContentType:   &i.File.MimeType,
	}, s3.WithPresignExpires(exp.Sub(time.Now())))
	if err != nil {
		return
	}

	i.S3UploadURL = &resp.URL
	i.S3UploadMethod = &resp.Method
	i.S3UploadExpiresAt = typex.P(exp.UnixMilli())
	return
}

type ImageFile struct {
	Size     int64  `dynamodbav:"size"`
	MimeType string `dynamodbav:"mime_type"`
}

func makeImage(file ImageFile) (res Image) {
	randomKey := uuid.New()
	res = Image{
		S3Bucket: imageS3Bucket,
		S3ObjectKey: path.Join(
			imageS3ObjectPrefix.Value(),
			sha3x.Sum256Base64(randomKey[:], base64.RawURLEncoding),
		),
		File: file,
	}
	res.genImageId()
	return
}

func imageId(str string) string {
	return sha3x.Sum256Base64([]byte(str), base64.RawURLEncoding)
}
