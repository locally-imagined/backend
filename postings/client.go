package postings

import (
	"backend/gen/postings"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type (
	// client is the postings service client interface
	Client interface {
		// Creates a post for the user and stores in s3 and postgres
		CreatePost(ctx context.Context, p *postings.CreatePostPayload) (*postings.CreatePostResult, error)

		// Gets first 25 posts in database
		GetPostPage(ctx context.Context, p *postings.GetPostPagePayload) (*postings.GetPostPageResult, error)

		// Gets first 25 posts for artist by id
		GetArtistPostPage(ctx context.Context, p *postings.GetArtistPostPagePayload) (*postings.GetArtistPostPageResult, error)

		// Gets first 25 filtered posts
		GetPostPageFiltered(ctx context.Context, p *postings.GetPostPageFilteredPayload) (*postings.GetPostPageFilteredResult, error)

		// Gets all image ids associated with a post
		GetImagesForPost(ctx context.Context, p *postings.GetImagesForPostPayload) (*postings.GetImagesForPostResult, error)

		// Deletes a post from database and images from s3
		DeletePost(ctx context.Context, p *postings.DeletePostPayload) error

		// Edits post with params given
		EditPost(ctx context.Context, p *postings.EditPostPayload) (*postings.EditPostResult, error)

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
		profpicID    string
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
	INSERTPOST   string = "INSERT INTO Posts (postid, userid, title, description, price, medium, deliverytype) Values ($1, $2, $3, $4, $5, $6, $7)"
	INSERTIMAGES string = "INSERT INTO Images Values ($1, $2, $3)"
	GETUSERNAME  string = "SELECT username FROM users WHERE userid=$1"
	GETPOSTPAGE  string = `SELECT p.postid, p.userid, u.username, u.profpicid, p.title, p.description, 
					p.price, p.medium, p.sold, p.uploaddate, p.deliverytype, i.imgid FROM posts AS p
					LEFT JOIN (SELECT imgid, postid FROM images WHERE index=0) AS i ON p.postid=i.postid 
					LEFT JOIN users AS u ON p.userid = u.userid
					ORDER BY p.uploaddate OFFSET $1 ROWS FETCH NEXT 25 ROWS ONLY`
	GETPOSTPAGEFILTERED string = `SELECT p.postid, p.userid, p.title, p.description, 
					p.price, p.medium, p.sold, p.uploaddate, i.imgid, u.username FROM posts AS p LEFT 
					JOIN images AS i ON p.postid=i.postid LEFT JOIN users AS u ON p.userid = u.userid WHERE (i.index=0) AND ((LOWER(p.title) LIKE $1) OR 
					(LOWER(p.description) LIKE $1)) AND (p.uploaddate >= $2) AND (p.uploaddate <= $3) AND (p.medium LIKE $4) 
					ORDER BY p.uploaddate OFFSET $5 ROWS FETCH NEXT 25 ROWS ONLY`
	GETPOSTPAGEFORARTIST string = `SELECT p.postid, p.userid, u.username, u.profpicid, p.title, 
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

func (c *client) CreatePost(ctx context.Context, p *postings.CreatePostPayload) (*postings.CreatePostResult, error) {
	awsBucketName, svc := c.getS3Session()
	postID := uuid.New().String()
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	now := time.Now().Format(time.RFC3339)
	_, err = dbPool.Query(INSERTPOST, postID, ctx.Value("UserID").(string), p.Post.Title, p.Post.Description, p.Post.Price, p.Post.Medium, p.Post.Deliverytype)
	if err != nil {
		return nil, err
	}
	imageIDs := make([]string, 0)
	for index, content := range p.Post.Content {
		imageID := uuid.New().String()
		imageIDs = append(imageIDs, imageID)
		reader := strings.NewReader(string(content))
		// put the object in the bucket
		err := putImageToS3(ctx, svc, awsBucketName, imageID, reader)
		if err != nil {
			return nil, err
		}
		_, err = dbPool.Query(INSERTIMAGES, imageID, postID, index)
		if err != nil {
			return nil, err
		}
	}
	posted := &postings.PostResponse{
		Title:        p.Post.Title,
		Description:  p.Post.Description,
		Price:        p.Post.Price,
		ImageIDs:     imageIDs,
		PostID:       postID,
		Medium:       p.Post.Medium,
		Sold:         false,
		UploadDate:   now[0:10],
		Deliverytype: p.Post.Deliverytype,
		UserID:       ctx.Value("UserID").(string),
	}
	res := &postings.CreatePostResult{
		Posted: posted,
	}
	return res, nil
}

func (c *client) GetPostPage(ctx context.Context, p *postings.GetPostPagePayload) (*postings.GetPostPageResult, error) {
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	offset := p.Page * IMAGESPERPAGE
	var rows *sql.Rows
	rows, err = dbPool.Query(GETPOSTPAGE, offset)
	if err != nil {
		return nil, err
	}
	res := make([]*postings.PostResponse, 0)
	for rows.Next() {
		var row post
		if err := rows.Scan(&row.postID, &row.userID, &row.username, &row.profpicID, &row.postTitle, &row.postDesc, &row.price, &row.medium, &row.sold, &row.uploadDate, &row.deliverytype, &row.imageID); err != nil {
			return nil, err
		}
		imageID := make([]string, 0)
		imageID = append(imageID, row.imageID)
		res = append(res, &postings.PostResponse{ProfpicID: row.profpicID, UserID: row.userID, Username: row.username, Title: row.postTitle, Description: row.postDesc, Price: row.price, ImageIDs: imageID, PostID: row.postID, UploadDate: row.uploadDate, Medium: row.medium, Sold: row.sold, Deliverytype: row.deliverytype})
	}
	return &postings.GetPostPageResult{Posts: res}, err
}

func (c *client) GetArtistPostPage(ctx context.Context, p *postings.GetArtistPostPagePayload) (*postings.GetArtistPostPageResult, error) {
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	offset := p.Page * IMAGESPERPAGE

	rows, err := dbPool.Query(GETPOSTPAGEFORARTIST, p.UserID, offset)
	if err != nil {
		return nil, err
	}
	res := make([]*postings.PostResponse, 0)
	for rows.Next() {
		var row post
		if err := rows.Scan(&row.postID, &row.userID, &row.username, &row.profpicID, &row.postTitle, &row.postDesc, &row.price, &row.medium, &row.sold, &row.uploadDate, &row.deliverytype, &row.imageID); err != nil {
			return nil, err
		}
		imageID := make([]string, 0)
		imageID = append(imageID, row.imageID)
		res = append(res, &postings.PostResponse{ProfpicID: row.profpicID, UserID: row.userID, Username: row.username, Title: row.postTitle, Description: row.postDesc, Price: row.price, ImageIDs: imageID, PostID: row.postID, UploadDate: row.uploadDate, Medium: row.medium, Sold: row.sold, Deliverytype: row.deliverytype})
	}
	return &postings.GetArtistPostPageResult{Posts: res}, err
}

func (c *client) GetPostPageFiltered(ctx context.Context, p *postings.GetPostPageFilteredPayload) (*postings.GetPostPageFilteredResult, error) {
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	offset := p.Page * IMAGESPERPAGE

	keyword := "%%"
	start := "2000-01-01"
	end := "2025-01-01"
	medium := "%%"
	if p.Keyword != nil {
		keyword = "%" + *p.Keyword + "%"
	}
	if p.StartDate != nil {
		start = *p.StartDate
	}
	if p.EndDate != nil {
		end = *p.EndDate
	}
	if p.Medium != nil {
		medium = "%" + *p.Medium + "%"
	}
	rows, err := dbPool.Query(GETPOSTPAGEFILTERED, keyword, start, end, medium, offset)
	if err != nil {
		return nil, err
	}
	res := make([]*postings.PostResponse, 0)
	for rows.Next() {
		var row post
		if err := rows.Scan(&row.postID, &row.userID, &row.postTitle, &row.postDesc, &row.price, &row.medium, &row.sold, &row.uploadDate, &row.imageID, &row.username); err != nil {
			return nil, err
		}
		imageID := make([]string, 0)
		imageID = append(imageID, row.imageID)
		res = append(res, &postings.PostResponse{UserID: row.userID, Title: row.postTitle, Description: row.postDesc, Price: row.price, ImageIDs: imageID, PostID: row.postID, UploadDate: row.uploadDate, Medium: row.medium, Sold: row.sold, Username: row.username})
	}
	return &postings.GetPostPageFilteredResult{Posts: res}, err
}

func (c *client) GetImagesForPost(ctx context.Context, p *postings.GetImagesForPostPayload) (*postings.GetImagesForPostResult, error) {
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	rows, err := dbPool.Query(SELECTIMAGES, p.PostID)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	for rows.Next() {
		var row image
		if err := rows.Scan(&row.imageID); err != nil {
			return nil, err
		}
		res = append(res, row.imageID)
	}
	return &postings.GetImagesForPostResult{Images: res}, err
}

func (c *client) DeletePost(ctx context.Context, p *postings.DeletePostPayload) error {
	dbPool, err := c.openDB()
	if err != nil {
		return err
	}
	defer dbPool.Close()
	rows, err := dbPool.Query(SELECTUSERID, p.PostID)
	if err != nil {
		return err
	}
	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return err
		}
	}
	if userID != ctx.Value("UserID").(string) {
		return err
	}

	rows, err = dbPool.Query(SELECTIMAGES, p.PostID)
	if err != nil {
		return err
	}
	imageIDs := make([]string, 0)
	for rows.Next() {
		var imageID string
		if err := rows.Scan(&imageID); err != nil {
			return err
		}
		imageIDs = append(imageIDs, imageID)
	}

	_, err = dbPool.Query(DELETEIMAGES, p.PostID)
	if err != nil {
		return err
	}
	_, err = dbPool.Query(DELETEPOST, p.PostID)
	if err != nil {
		return err
	}
	awsBucketName, svc := c.getS3Session()
	// delete the images in the bucket
	for _, imageID := range imageIDs {
		err = deleteImageFromS3(ctx, svc, awsBucketName, imageID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) EditPost(ctx context.Context, p *postings.EditPostPayload) (*postings.EditPostResult, error) {
	dbPool, err := c.openDB()
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	rows, err := dbPool.Query(SELECTUSERID, p.PostID)
	if err != nil {
		return nil, err
	}
	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
	}
	if userID != ctx.Value("UserID").(string) {
		return nil, err
	}
	query1 := "UPDATE posts SET "
	query2 := "=$1 WHERE postID=$2"
	if p.Title != nil {
		query := query1 + "title" + query2
		_, err = dbPool.Query(query, *p.Title, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Description != nil {
		query := query1 + "description" + query2
		_, err = dbPool.Query(query, *p.Description, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Price != nil {
		query := query1 + "price" + query2
		_, err = dbPool.Query(query, *p.Price, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Medium != nil {
		query := query1 + "medium" + query2
		_, err = dbPool.Query(query, *p.Medium, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Sold != nil {
		query := query1 + "sold" + query2
		_, err = dbPool.Query(query, *p.Sold, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	if p.Deliverytype != nil {
		query := query1 + "deliverytype" + query2
		_, err = dbPool.Query(query, *p.Deliverytype, p.PostID)
		if err != nil {
			return nil, err
		}
	}

	// what if image is the last pic? dont let user delete last pic?
	if p.ImageID != nil {
		_, err = dbPool.Query(UPDATEINDEX, p.PostID, *p.ImageID)
		if err != nil {
			return nil, err
		}
		_, err = dbPool.Query(DELETEIMAGE, *p.ImageID)
		if err != nil {
			return nil, err
		}
		awsBucketName, svc := c.getS3Session()
		err = deleteImageFromS3(ctx, svc, awsBucketName, *p.ImageID)
		if err != nil {
			return nil, err
		}
	}
	if p.Content.Content != nil {
		fmt.Printf("%v", p.Content)
		fmt.Printf("%s", *p.Content.Content)
		imageID := uuid.New().String()
		awsBucketName, svc := c.getS3Session()
		reader := strings.NewReader(string(*p.Content.Content))
		// put the object in the bucket
		err := putImageToS3(ctx, svc, awsBucketName, imageID, reader)
		if err != nil {
			return nil, err
		}
		_, err = dbPool.Query(ADDIMAGE, imageID, p.PostID, p.PostID)
		if err != nil {
			return nil, err
		}
	}
	rows, err = dbPool.Query(GETEDITEDINFO, p.PostID)
	if err != nil {
		return nil, err
	}
	var row post
	for rows.Next() {
		if err := rows.Scan(&row.postTitle, &row.postDesc, &row.price, &row.medium, &row.sold, &row.deliverytype, &row.uploadDate); err != nil {
			return nil, err
		}
	}
	rows, err = dbPool.Query(SELECTIMAGES, p.PostID)
	if err != nil {
		return nil, err
	}
	imageIDs := make([]string, 0)
	for rows.Next() {
		var imageID string
		if err := rows.Scan(&imageID); err != nil {
			return nil, err
		}
		imageIDs = append(imageIDs, imageID)
	}
	posted := &postings.PostResponse{
		Title:        row.postTitle,
		Description:  row.postDesc,
		Price:        row.price,
		ImageIDs:     imageIDs,
		PostID:       p.PostID,
		Medium:       row.medium,
		Sold:         row.sold,
		UploadDate:   row.uploadDate,
		Deliverytype: row.deliverytype,
		UserID:       ctx.Value("UserID").(string),
	}
	res := &postings.EditPostResult{
		Posted: posted,
	}
	return res, nil

}
