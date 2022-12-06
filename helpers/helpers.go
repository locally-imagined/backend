package helpers

import (
	"context"
	"database/sql"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func OpenDB(dbURL string) (*sql.DB, error) {
	return sql.Open("postgres", dbURL)
}

func GetS3Session(awsAccessKey, awsSecretKey, awsRegion, awsBucketName string) (string, *s3.S3) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	}))
	svc := s3.New(sess)
	return awsBucketName, svc
}

func PutImageToS3(ctx context.Context, svc *s3.S3, awsBucketName, imageID string, reader *strings.Reader) error {
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String("public/" + imageID),
		Body:   reader,
	})
	return err
}

func DeleteImageFromS3(ctx context.Context, svc *s3.S3, awsBucketName, imageID string) error {
	_, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String("public/" + imageID),
	})
	return err
}
