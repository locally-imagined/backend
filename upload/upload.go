package upload

import (
	"backend/gen/upload"
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Service struct{}

func (s *Service) UploadPhoto(ctx context.Context, p *upload.UploadPhotoPayload) (*upload.UploadPhotoResult, error) {
	uploadBool := false
	star := "*"
	// get info from os variables
	aws_access_key := os.Getenv("BUCKETEER_AWS_ACCESS_KEY_ID")
	aws_secret_key := os.Getenv("BUCKETEER_AWS_SECRET_ACCESS_KEY")
	aws_region := os.Getenv("BUCKETEER_AWS_REGION")
	aws_bucket_name := os.Getenv("BUCKETEER_BUCKET_NAME")

	// create new s3 session
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(aws_region),
		Credentials: credentials.NewStaticCredentials(aws_access_key, aws_secret_key, ""),
	}))
	// create a new instance of the service's client iwth a session.
	svc := s3.New(sess)
	// put the bytes into a reader, bytes must be in base 64 for this to work
	reader := strings.NewReader(string(p.Content))

	new_uuid := uuid.NewString()
	// put the object in the bucket
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(aws_bucket_name),
		Key:    aws.String(new_uuid),
		Body:   reader,
	})
	if err != nil {
		uploadBool = false
	}
	// insert uuid into the database
	return &upload.UploadPhotoResult{Success: &uploadBool, AccessControlAllowOrigin: &star}, nil
}
