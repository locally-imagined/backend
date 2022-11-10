package upload

// should be package postings
import (
	"backend/auth"
	"backend/gen/postings"
	"context"
	"database/sql"
	"log"
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
	ErrUnauthorized error = postings.Unauthorized("invalid jwt")
)

func (s *Service) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	tok := auth.DecodeToken(token)
	if tok == nil {
		return ctx, ErrUnauthorized
	}
	// 3. add authInfo to context

	claims := tok.Claims.(jwt.MapClaims)
	ctx = context.WithValue(ctx, "Username", claims["Username"])
	ctx = context.WithValue(ctx, "UserID", claims["UserID"])
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

func (s *Service) CreatePost(ctx context.Context, p *postings.CreatePostPayload) (*postings.CreatePostResult, error) {
	// this is really CreatePost now
	// create a different endpoint UploadPhoto that takes in a postID
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
	reader := strings.NewReader(string(p.Post.Content))

	imageID := uuid.New().String()
	// put the object in the bucket
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String("public/" + imageID),
		Body:   reader,
	})
	if err != nil {
		return nil, err
	}

	postID := uuid.New().String()
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339)
	_, err = dbPool.Query("INSERT INTO Posts Values ($1, $2, $3, $4, $5, $6)", postID, ctx.Value("UserID").(string), p.Post.Title, p.Post.Description, p.Post.Price, now)
	if err != nil {
		return nil, err
	}
	_, err = dbPool.Query("INSERT INTO Images Values ($1, $2, $3)", imageID, postID, 0)
	if err != nil {
		return nil, err
	}
	posted := &postings.PostResponse{
		Title:       p.Post.Title,
		Description: p.Post.Description,
		Price:       p.Post.Price,
		ImageID:     imageID,
		PostID:      postID,
		UploadDate:  now,
	}
	res := &postings.CreatePostResult{
		Posted: posted,
	}
	return res, nil
}

type image struct {
	imageID string
	index   int
}

type post struct {
	postID     string
	userID     string
	postTitle  string
	postDesc   string
	price      string
	uploadDate string
	imageID    string
}

func (s *Service) GetPostPage(ctx context.Context, p *postings.GetPostPagePayload) (*postings.GetPostPageResult, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	// queryString := `SELECT p.postid, p.userid, p.title,
	// p.description, p.price, p.uploaddate, i.imgid
	// FROM posts AS p LEFT JOIN images
	// AS i ON p.postid = i.postid WHERE i.index =
	// 0 OFFSET $1 ROWS FETCH NEXT 25 ROWS ONLY`
	offset := p.Page * 25
	rows, err := dbPool.Query("SELECT p.postid, p.userid, p.title, p.description, p.price, p.uploaddate, i.imgid FROM posts AS p LEFT JOIN images AS i ON p.postid=i.postid WHERE i.index=0 OFFSET $1 ROWS FETCH NEXT 25 ROWS ONLY", offset)

	defer dbPool.Close()

	if err == sql.ErrNoRows {
		return nil, err
	}

	res := make([]*postings.PostResponse, 0)
	i := 0
	for rows.Next() {
		var row post
		if err := rows.Scan(&row.postID, &row.userID, &row.postTitle, &row.postDesc, &row.price, &row.uploadDate, &row.imageID); err != nil {
			log.Fatal(err)
		}
		res = append(res, &postings.PostResponse{Title: row.postTitle, Description: row.postDesc, Price: row.price, ImageID: row.imageID, PostID: row.postID, UploadDate: row.uploadDate})
		i = i + 1
	}
	return &postings.GetPostPageResult{Posts: res}, err
}
