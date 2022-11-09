package upload

import (
	"backend/auth"
	"backend/gen/upload"
	"context"
	"time"

	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"goa.design/goa/v3/security"
)

type Service struct{}

var (
	// ErrUnauthorized is the error returned by Login when the request credentials
	// are invalid.
	ErrUnauthorized error = upload.Unauthorized("invalid jwt")
)

func (s *Service) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	tok := auth.DecodeToken(token)
	if tok == nil {
		return ctx, ErrUnauthorized
	}
	// 3. add authInfo to context

	claims := tok.Claims.(jwt.MapClaims)
	ctx = context.WithValue(ctx, "Name", claims["Name"])
	var exp time.Time
	var now time.Time = time.Now()
	switch iat := claims["iat"].(type) {
	case float64:
		exp = time.Unix(int64(iat), 0)
	}
	if exp.Add(time.Hour * 2).Before(now) {
		return ctx, ErrUnauthorized
	}

	return ctx, nil
}

func (s *Service) UploadPhoto(ctx context.Context, p *upload.UploadPhotoPayload) (*upload.UploadPhotoResult, error) {
	star := "*"
	// get info from os variables
	awsAccessKey := os.Getenv("BUCKETEER_AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("BUCKETEER_AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("BUCKETEER_AWS_REGION")
	awsBucketName := os.Getenv("BUCKETEER_BUCKET_NAME")

	// create new s3 session
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	}))
	// create a new instance of the service's client iwth a session.
	svc := s3.New(sess)
	// put the bytes into a reader, bytes must be in base 64 for this to work
	reader := strings.NewReader(string(p.Content))

	postID := uuid.NewString()
	// put the object in the bucket
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String("public/" + postID),
		Body:   reader,
	})
	if err != nil {
		return &upload.UploadPhotoResult{Success: &star}, err
	}
	// x := resp.
	// 	// insert uuid into the database
	// 	fmt.Printf("response %s", awsutil.Prettify(resp))
	// params := &s3.ListObjectsInput{
	// 	Bucket: aws.String("bucket"),
	// }
	return &upload.UploadPhotoResult{Success: &star}, nil
}
