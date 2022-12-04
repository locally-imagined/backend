package users

import (
	"backend/gen/users"
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/lib/pq"
)

type (
	// client is the users service client interface
	Client interface {
		// Updates user bio for given userid
		UpdateBio(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error)

		// Retrieves contact info for given userid
		GetContactInfo(ctx context.Context, p *users.GetContactInfoPayload) (*users.GetContactInfoResult, error)

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

	image struct {
		imageID string
	}

	post struct {
		postID       string
		userID       string
		postTitle    string
		postDesc     string
		price        string
		uploadDate   string
		imageID      string
		medium       string
		sold         bool
		deliverytype string
		username     string
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
	INSERTPOST   string = "INSERT INTO Posts (postid, userid, title, description, price, medium, deliverytype) Values ($1, $2, $3, $4, $5, $6, $7)"
	INSERTIMAGES string = "INSERT INTO Images Values ($1, $2, $3)"
	GETUSERNAME  string = "SELECT username FROM users WHERE userid=$1"
	GETPOSTPAGE  string = `SELECT p.postid, p.userid, u.username, p.title, p.description, 
					p.price, p.medium, p.sold, p.uploaddate, p.deliverytype, i.imgid FROM posts AS p
					LEFT JOIN (SELECT imgid, postid FROM images WHERE index=0) AS i ON p.postid=i.postid 
					LEFT JOIN users AS u ON p.userid = u.userid
					ORDER BY p.uploaddate OFFSET $1 ROWS FETCH NEXT 25 ROWS ONLY`
	GETPOSTPAGEWITHKEYWORD string = `SELECT p.postid, p.userid, p.title, p.description, 
					p.price, p.medium, p.sold, p.uploaddate, p.deliverytype, i.imgid FROM posts AS p LEFT 
					JOIN images AS i ON p.postid = i.postid WHERE i.index=0 AND 
					((LOWER(p.title) LIKE $1) OR (LOWER(p.description) LIKE $2))
					ORDER BY p.uploaddate OFFSET $3 ROWS FETCH NEXT 25 ROWS ONLY`
	GETPOSTPAGEFORARTIST string = `SELECT p.postid, p.userid, u.username, p.title, 
					p.description, p.price, p.medium, p.sold, p.uploaddate, p.deliverytype, i.imgid FROM posts AS p
					LEFT JOIN (SELECT imgid, postid FROM images WHERE index=0) AS i ON p.postid=i.postid 
					LEFT JOIN users AS u ON p.userid = u.userid where p.userid=$1      
					ORDER BY p.uploaddate OFFSET $2 ROWS FETCH NEXT 25 ROWS ONLY;`
	SELECTIMAGES  string = "SELECT imgid from images where postid=$1 ORDER BY index"
	SELECTUSERID  string = "SELECT userID from Posts where postID=$1"
	DELETEIMAGES  string = "DELETE FROM images WHERE postID=$1"
	DELETEPOST    string = "DELETE FROM posts WHERE postID=$1"
	DELETEIMAGE   string = "DELETE FROM images where imgid=$1"
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
	_, err = dbPool.Query(UPDATEBIO, *p.Bio, ctx.Value("UserID").(string))
	if err != nil {
		return nil, err
	}

	rows, err := dbPool.Query(GETUSER, ctx.Value("UserID").(string))
	var resp users.UpdateBioResult
	for rows.Next() {
		var row users.User
		if err := rows.Scan(&row.FirstName, &row.LastName, &row.Phone, &row.Email); err != nil {
			log.Fatal(err)
			return nil, err
		}
		resp = users.UpdateBioResult{UpdatedUser: &row}
		//res = append(res, &users.User{FirstName: row.firstname, LastName: row.lastname, Phone: row.phone, Email: row.email})
	}
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

	var resp users.GetContactInfoResult
	for rows.Next() {
		var row users.User
		if err := rows.Scan(&row.FirstName, &row.LastName, &row.Phone, &row.Email); err != nil {
			log.Fatal(err)
			return nil, err
		}
		resp = users.GetContactInfoResult{ContactInfo: &row}
		//res = append(res, &users.User{FirstName: row.firstname, LastName: row.lastname, Phone: row.phone, Email: row.email})
	}
	return &resp, err
}
