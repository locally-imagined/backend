package users

import (
	"backend/gen/users"
	"context"
	"database/sql"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type (
	// client is the users service client interface
	Client interface {
		// Updates user bio for given userid
		UpdateBio(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error)

		// Retrieves contact info for given userid
		GetContactInfo(ctx context.Context, p *users.GetContactInfoPayload) (*users.GetContactInfoResult, error)

		// Updates user profile picture
		UpdateProfilePicture(ctx context.Context, p *users.UpdateProfilePicturePayload) (*users.UpdateProfilePictureResult, error)

		// Opens postgres db connection
		openDB() (*sql.DB, error)

		// Creates s3 session
		getS3Session() (string, *s3.S3)
	}

	client struct {
		awsAccessKey  string
		awsSecretKey  string
		awsRegion     string
		awsBucketName string
		dbURL         string
	}
)

func New(awsAccessKey, awsSecretKey, awsRegion, awsBucketName, dbURL string) Client {
	return &client{
		awsAccessKey:  awsAccessKey,
		awsSecretKey:  awsSecretKey,
		awsRegion:     awsRegion,
		awsBucketName: awsBucketName,
		dbURL:         dbURL,
	}
}

var (
	// ErrUnauthorized is the error returned by Login when the request credentials
	// are invalid.

	// add prepares before queries?
	//
	GETPROFILEPIC    string = "SELECT profpicid FROM users WHERE userid=$1 AND profpicid IS NOT NULL"
	UPDATEPROFILEPIC string = "UPDATE users SET profpicid=$1 WHERE userid=$2"

	UPDATEINDEX   string = "UPDATE images SET index = index - 1 WHERE (postid=$1 AND index>(SELECT index FROM images WHERE imgid=$2))"
	ADDIMAGE      string = "INSERT INTO images VALUES($1, $2, (SELECT MAX(index) FROM images where postID=$3) + 1)"
	GETEDITEDINFO string = "SELECT title, description, price, medium, sold, deliverytype, uploaddate FROM posts where postID=$1"
	IMAGESPERPAGE int    = 25
	UPDATEBIO     string = "UPDATE users SET bio = $1 WHERE userID = $2"
	GETUSER       string = "SELECT firstname, lastname, phone, email FROM users WHERE userID = $1"
)

func (c *client) openDB() (*sql.DB, error) {
	return sql.Open("postgres", c.dbURL)
}

func (c *client) getS3Session() (string, *s3.S3) {
	awsAccessKey := c.awsAccessKey
	awsSecretKey := c.awsSecretKey
	awsRegion := c.awsRegion
	awsBucketName := c.awsBucketName
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	}))
	// create a new instance of the service's client with a session.
	svc := s3.New(sess)
	return awsBucketName, svc
}

func putImageToS3(ctx context.Context, svc *s3.S3, awsBucketName, imageID string, reader *strings.Reader) error {
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String("public/" + imageID),
		Body:   reader,
	})
	return err
}

func deleteImageFromS3(ctx context.Context, svc *s3.S3, awsBucketName, imageID string) error {
	_, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String("public/" + imageID),
	})
	return err
}

func (c *client) UpdateBio(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error) {
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	_, err = dbPool.Query(UPDATEBIO, *p.Bio.Bio, ctx.Value("UserID").(string))
	if err != nil {
		return nil, err
	}
	rows, err := dbPool.Query(GETUSER, ctx.Value("UserID").(string))
	var row users.User
	for rows.Next() {
		if err := rows.Scan(&row.FirstName, &row.LastName, &row.Phone, &row.Email); err != nil {
			return nil, err
		}
	}
	resp := users.UpdateBioResult{UpdatedUser: &row}
	return &resp, err
}

func (c *client) GetContactInfo(ctx context.Context, p *users.GetContactInfoPayload) (*users.GetContactInfoResult, error) {
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()

	var rows *sql.Rows
	rows, err = dbPool.Query(GETUSER, p.UserID)
	if err != nil {
		return nil, err
	}

	var row users.User
	for rows.Next() {
		if err := rows.Scan(&row.FirstName, &row.LastName, &row.Phone, &row.Email); err != nil {
			return nil, err
		}
	}
	resp := users.GetContactInfoResult{ContactInfo: &row}
	return &resp, err
}

func (c *client) UpdateProfilePicture(ctx context.Context, p *users.UpdateProfilePicturePayload) (*users.UpdateProfilePictureResult, error) {
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()

	var rows *sql.Rows
	rows, err = dbPool.Query(GETPROFILEPIC, ctx.Value("UserID").(string))
	if (err != nil) && (err != sql.ErrNoRows) {
		return nil, err
	}
	bucketName, svc := c.getS3Session()
	if err != sql.ErrNoRows {
		var oldID string
		for rows.Next() {
			if err := rows.Scan(&oldID); err != nil {
				return nil, err
			}
		}
		err = deleteImageFromS3(ctx, svc, bucketName, oldID)
		if err != nil {
			return nil, err
		}
	}
	newID := uuid.New().String()
	_, err = dbPool.Query(UPDATEPROFILEPIC, newID, ctx.Value("UserID").(string))
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(*p.Content.Content))
	err = putImageToS3(ctx, svc, bucketName, newID, reader)
	if err != nil {
		return nil, err
	}
	photo := users.ProfilePhoto{ImageID: &newID}
	resp := users.UpdateProfilePictureResult{ImageID: &photo}
	return &resp, err
}
