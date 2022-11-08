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
	row, err := dbPool.Query("SELECT password from test_users where username=$1", user)
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
	return ctx, nil
}

func (s *Service) Login(ctx context.Context, p *login.LoginPayload) (*login.LoginResult, error) {
	access := "*"
	creds := "true"
	token, err := auth.MakeToken(p.Username)
	if err != nil {
		return &login.LoginResult{JWT: nil, AccessControlAllowOrigin: &access}, err
	}
	return &login.LoginResult{JWT: &token, AccessControlAllowOrigin: &access, AccessControlAllowCredentials: &creds}, nil
}
