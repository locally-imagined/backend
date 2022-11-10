package signup

import (
	"backend/auth"
	"backend/gen/login"
	"backend/gen/signup"
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/google/uuid"
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
	hashedPassword := auth.ShaHashing(pass)
	var value string = ""
	// Query for a value based on a single row.
	row, err := dbPool.Query("SELECT username from test_users where username=$1", user)
	if err != nil {
		return ctx, ErrUnauthorized
	}
	for row.Next() {
		if err := row.Scan(&value); err != nil {
			log.Fatal(err)
		}
	}
	if value != "" {
		return ctx, ErrUnauthorized
	}

	// double check this
	userID := uuid.New()

	_, err = dbPool.Query("INSERT INTO Users (username, password) Values ($1, $2)", user, hashedPassword)
	// insert other items including generated userID
	if err != nil {
		return ctx, ErrUnauthorized
	}
	ctx = context.WithValue(ctx, "UserID", userID)
	return ctx, nil
}

func (s *Service) Signup(ctx context.Context, p *signup.SignupPayload) (*signup.SignupResult, error) {
	token, err := auth.MakeToken(p.Username, ctx.Value("UserID").(string))
	if err != nil {
		return &signup.SignupResult{JWT: nil}, err
	}
	return &signup.SignupResult{JWT: &token}, nil
}
