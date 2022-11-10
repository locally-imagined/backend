package upload

// should be package postings
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
	// this is really CreatePost now
	// create a different endpoint UploadPhoto that takes in a postID
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
	// postID = new uuid
	// insert into 'posts' postID, p.title, p.desc, p.price, currentDate
	// insert into 'images', imageID (maybe), postID, https://bucketeer-8e1fe0c2-5dfb-4787-8878-55a22a5940a8.s3.amazonaws.com/public/postID, isThumbnail=tru
	return &upload.UploadPhotoResult{Success: &star}, nil
}
