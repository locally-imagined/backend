package login

import (
	"backend/auth"
	"backend/gen/login"
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"goa.design/goa/v3/security"
)

type Service struct{}

var (
	// ErrUnauthorized is the error returned by Login when the request credentials
	// are invalid.
	ErrUnauthorized error = login.Unauthorized("invalid username and password combination")
)

func (s *Service) BasicAuth(ctx context.Context, user, pass string, scheme *security.BasicScheme) (context.Context, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return ctx, ErrUnauthorized
	}
	defer dbPool.Close()
	var password string
	hashedPassword := auth.ShaHashing(pass)
	row, err := dbPool.Query("SELECT passhash from users where username=$1", user)
	if err == sql.ErrNoRows {
		return ctx, ErrUnauthorized
	}
	for row.Next() {
		if err := row.Scan(&password); err != nil {
			log.Fatal(err)
		}
	}
	if hashedPassword != password {
		return ctx, ErrUnauthorized
	}
	var userID string
	row, err = dbPool.Query("SELECT userID from users where username=$1", user)
	if err == sql.ErrNoRows {
		return ctx, ErrUnauthorized
	}
	for row.Next() {
		if err := row.Scan(&userID); err != nil {
			log.Fatal(err)
		}
	}
	// add userID into context
	ctx = context.WithValue(ctx, "UserID", userID)
	return ctx, nil
}

func (s *Service) Login(ctx context.Context, p *login.LoginPayload) (*login.LoginResult, error) {
	// add userID into token
	UserID := ctx.Value("UserID").(string)
	token, err := auth.MakeToken(p.Username, ctx.Value("UserID").(string))
	if err != nil {
		return nil, err
	}
	resp := &login.LoginResponse{
		UserID: &UserID,
		JWT:    &token,
	}
	res := &login.LoginResult{
		LoginResponse: resp,
	}
	return res, nil
}
