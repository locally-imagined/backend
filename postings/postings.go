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

	postID := uuid.New().String()
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339)
	_, err = dbPool.Query("INSERT INTO Posts Values ($1, $2, $3, $4, $5, $6)", postID, ctx.Value("UserID").(string), p.Post.Title, p.Post.Description, p.Post.Price, p.Post.Medium)
	if err != nil {
		return nil, err
	}
	imageIDs := make([]string, 0)
	for index, content := range p.Post.Content {

		imageID := uuid.New().String()
		imageIDs = append(imageIDs, imageID)
		reader := strings.NewReader(string(content))
		// put the object in the bucket
		_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(awsBucketName),
			Key:    aws.String("public/" + imageID),
			Body:   reader,
		})
		if err != nil {
			return nil, err
		}
		_, err = dbPool.Query("INSERT INTO Images Values ($1, $2, $3)", imageID, postID, index)
		if err != nil {
			return nil, err
		}
	}

	posted := &postings.PostResponse{
		Title:       p.Post.Title,
		Description: p.Post.Description,
		Price:       p.Post.Price,
		ImageIDs:    imageIDs,
		PostID:      postID,
		Medium:      p.Post.Medium,
		Sold:        false,
		UploadDate:  now,
	}
	res := &postings.CreatePostResult{
		Posted: posted,
	}
	return res, nil
}

type image struct {
	imageID string
}

type post struct {
	postID     string
	userID     string
	postTitle  string
	postDesc   string
	price      string
	uploadDate string
	imageID    string
	medium     string
	sold       bool
}

func (s *Service) GetPostPage(ctx context.Context, p *postings.GetPostPagePayload) (*postings.GetPostPageResult, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	offset := p.Page * 25
	rows, err := dbPool.Query("SELECT p.postid, p.userid, p.title, p.description, p.price, p.medium, p.sold, p.uploaddate, i.imgid FROM posts AS p LEFT JOIN images AS i ON p.postid=i.postid WHERE i.index=0 ORDER BY p.uploaddate OFFSET $1 ROWS FETCH NEXT 25 ROWS ONLY", offset)

	defer dbPool.Close()

	if err == sql.ErrNoRows {
		return nil, err
	}

	res := make([]*postings.PostResponse, 0)
	for rows.Next() {
		var row post
		if err := rows.Scan(&row.postID, &row.userID, &row.postTitle, &row.postDesc, &row.price, &row.medium, &row.sold, &row.uploadDate, &row.imageID); err != nil {
			log.Fatal(err)
			return nil, err
		}
		imageID := make([]string, 0)
		imageID = append(imageID, row.imageID)
		res = append(res, &postings.PostResponse{Title: row.postTitle, Description: row.postDesc, Price: row.price, ImageIDs: imageID, PostID: row.postID, UploadDate: row.uploadDate, Medium: row.medium, Sold: row.sold})
	}
	return &postings.GetPostPageResult{Posts: res}, err
}

func (s *Service) GetImagesForPost(ctx context.Context, p *postings.GetImagesForPostPayload) (*postings.GetImagesForPostResult, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	rows, err := dbPool.Query("SELECT imgid from images where postid=$1 ORDER BY index", p.PostID)

	defer dbPool.Close()

	if err == sql.ErrNoRows {
		return nil, err
	}

	res := make([]string, 0)
	for rows.Next() {
		var row image
		if err := rows.Scan(&row.imageID); err != nil {
			log.Fatal(err)
		}
		res = append(res, row.imageID)
	}
	return &postings.GetImagesForPostResult{Images: res}, err
}

func (s *Service) DeletePost(ctx context.Context, p *postings.DeletePostPayload) error {

	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	rows, err := dbPool.Query("SELECT userID from Posts where postID=$1", p.PostID)
	if err != nil {
		return err
	}
	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			log.Fatal(err)
			return err
		}
	}
	if userID != ctx.Value("UserID").(string) {
		return err
	}

	rows, err = dbPool.Query("SELECT imgid from images where postID=$1", p.PostID)
	if err != nil {
		return err
	}
	imageIDs := make([]string, 0)
	for rows.Next() {
		var imageID string
		if err := rows.Scan(&imageID); err != nil {
			log.Fatal(err)
			return err
		}
		imageIDs = append(imageIDs, imageID)
	}

	_, err = dbPool.Query("DELETE FROM images WHERE postID=$1", p.PostID)
	if err != nil {
		return err
	}
	_, err = dbPool.Query("DELETE FROM posts WHERE postID=$1", p.PostID)
	if err != nil {
		return err
	}
	awsAccessKey := os.Getenv("BUCKETEER_AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("BUCKETEER_AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("BUCKETEER_AWS_REGION")
	awsBucketName := os.Getenv("BUCKETEER_BUCKET_NAME")
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	}))
	svc := s3.New(sess)

	// delete the images in the bucket
	for _, imageID := range imageIDs {
		_, err = svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(awsBucketName),
			Key:    aws.String("public/" + imageID),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) EditPost(ctx context.Context, p *postings.EditPostPayload) (*postings.EditPostResult, error) {

	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	rows, err := dbPool.Query("SELECT userID from Posts where postID=$1", p.PostID)
	if err != nil {
		return nil, err
	}
	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	if userID != ctx.Value("UserID").(string) {
		return nil, err
	}

	if p.Title != nil {
		_, err = dbPool.Query("UPDATE posts SET title=$1 WHERE postID=$2", *p.Title, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Description != nil {
		_, err = dbPool.Query("UPDATE posts SET description=$1 WHERE postID=$2", *p.Description, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Price != nil {
		_, err = dbPool.Query("UPDATE posts SET price=$1 WHERE postID=$2", *p.Price, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Medium != nil {
		_, err = dbPool.Query("UPDATE posts SET medium=$1 WHERE postID=$2", *p.Medium, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Sold != nil {
		_, err = dbPool.Query("UPDATE posts SET sold=$1 WHERE postID=$2", *p.Sold, p.PostID)
		if err != nil {
			return nil, err
		}
	}

	// what if image is the last pic? dont let user delete last pic?
	if p.ImageID != nil {
		_, err = dbPool.Query("UPDATE images SET index = index - 1 WHERE (postid=$1 AND index>(SELECT index FROM images WHERE imgid=$2))", p.PostID, *p.ImageID)
		if err != nil {
			return nil, err
		}
		_, err = dbPool.Query("DELETE FROM where imgid=$1", *p.ImageID)
		if err != nil {
			return nil, err
		}
	}
	if p.Content != nil {
		imageID := uuid.New().String()
		awsAccessKey := os.Getenv("BUCKETEER_AWS_ACCESS_KEY_ID")
		awsSecretKey := os.Getenv("BUCKETEER_AWS_SECRET_ACCESS_KEY")
		awsRegion := os.Getenv("BUCKETEER_AWS_REGION")
		awsBucketName := os.Getenv("BUCKETEER_BUCKET_NAME")
		sess := session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
		}))
		svc := s3.New(sess)
		reader := strings.NewReader(string(*p.Content))
		// put the object in the bucket
		_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(awsBucketName),
			Key:    aws.String("public/" + imageID),
			Body:   reader,
		})
		if err != nil {
			return nil, err
		}
		_, err = dbPool.Query("INSERT INTO images VALUES($1, $2, (SELECT MAX(index) FROM images where postID=$3) + 1)", imageID, p.PostID, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	rows, err = dbPool.Query("SELECT title, description, price, medium, sold, uploaddate FROM posts where postID=$1", p.PostID)
	if err != nil {
		return nil, err
	}
	var row post
	for rows.Next() {
		if err := rows.Scan(&row.postTitle, &row.postDesc, &row.price, &row.medium, &row.sold, &row.uploadDate); err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	rows, err = dbPool.Query("SELECT imgid from images where postID=$1", p.PostID)
	if err != nil {
		return nil, err
	}
	imageIDs := make([]string, 0)
	for rows.Next() {
		var imageID string
		if err := rows.Scan(&imageID); err != nil {
			log.Fatal(err)
			return nil, err
		}
		imageIDs = append(imageIDs, imageID)
	}
	posted := &postings.PostResponse{
		Title:       row.postTitle,
		Description: row.postDesc,
		Price:       row.price,
		ImageIDs:    imageIDs,
		PostID:      p.PostID,
		Medium:      row.medium,
		Sold:        row.sold,
		UploadDate:  row.uploadDate,
	}
	res := &postings.EditPostResult{
		Posted: posted,
	}
	return res, nil

}
