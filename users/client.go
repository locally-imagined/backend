package users

import (
	"backend/gen/users"
	"backend/helpers"
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type (
	// client is the users service client interface
	Client interface {
		// Updates user bio for given userid
		UpdateBio(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error)

		// Retrieves contact info for given userid
		GetUserInfo(ctx context.Context, p *users.GetUserInfoPayload) (*users.GetUserInfoResult, error)

		// Updates user profile picture
		UpdateProfilePicture(ctx context.Context, p *users.UpdateProfilePicturePayload) (*users.UpdateProfilePictureResult, error)
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
	GETPROFILEPIC    string = "SELECT profpicid FROM users WHERE userid=$1"
	UPDATEPROFILEPIC string = "UPDATE users SET profpicid=$1 WHERE userid=$2"
	UPDATEBIO        string = "UPDATE users SET bio = $1 WHERE userID = $2"
	GETUSER          string = "SELECT firstname, lastname, phone, email, bio, profpicid FROM users WHERE userID = $1"
)

func (c *client) UpdateBio(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error) {
	dbPool, err := helpers.OpenDB(c.dbURL)
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	_, err = dbPool.Query(UPDATEBIO, *p.Bio.Bio, ctx.Value("UserID").(string))
	if err != nil {
		return nil, err
	}
	rows, err := dbPool.Query(GETUSER, ctx.Value("UserID").(string))
	var user users.User
	for rows.Next() {
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.Phone, &user.Email, &user.Bio, &user.ProfpicID); err != nil {
			return nil, err
		}
	}
	resp := users.UpdateBioResult{User: &user}
	return &resp, err
}

func (c *client) GetUserInfo(ctx context.Context, p *users.GetUserInfoPayload) (*users.GetUserInfoResult, error) {
	dbPool, err := helpers.OpenDB(c.dbURL)
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()

	var rows *sql.Rows
	rows, err = dbPool.Query(GETUSER, p.UserID)
	if err != nil {
		return nil, err
	}

	var user users.User
	for rows.Next() {
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.Phone, &user.Email, &user.Bio, &user.ProfpicID); err != nil {
			return nil, err
		}
	}
	resp := users.GetUserInfoResult{User: &user}
	return &resp, err
}

func (c *client) UpdateProfilePicture(ctx context.Context, p *users.UpdateProfilePicturePayload) (*users.UpdateProfilePictureResult, error) {
	dbPool, err := helpers.OpenDB(c.dbURL)
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()

	var rows *sql.Rows
	rows, err = dbPool.Query(GETPROFILEPIC, ctx.Value("UserID").(string))
	if err != nil {
		return nil, err
	}
	bucketName, svc := helpers.GetS3Session(c.awsAccessKey, c.awsSecretKey, c.awsRegion, c.awsBucketName)
	var oldID string
	for rows.Next() {
		if err := rows.Scan(&oldID); err != nil {
			return nil, err
		}
	}
	if oldID != uuid.Nil.String() {
		err = helpers.DeleteImageFromS3(ctx, svc, bucketName, oldID)
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
	err = helpers.PutImageToS3(ctx, svc, bucketName, newID, reader)
	if err != nil {
		return nil, err
	}
	photo := users.ProfilePhoto{ImageID: &newID}
	resp := users.UpdateProfilePictureResult{ImageID: &photo}
	return &resp, err
}
