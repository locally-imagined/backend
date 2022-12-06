package signup

import (
	"backend/auth"
	"backend/gen/signup"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

type Service struct{}

var (
	// ErrUnauthorized is the error returned by signup when the request credentials
	// are invalid or if anything else goes wrong.
	ErrUnauthorized *goa.ServiceError = signup.MakeUnauthorized(fmt.Errorf("Unauthorized"))
)

func (s *Service) BasicAuth(ctx context.Context, user, pass string, scheme *security.BasicScheme) (context.Context, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return ctx, ErrUnauthorized
	}
	defer dbPool.Close()
	var value string = ""
	// Query for a value based on a single row.
	row, err := dbPool.Query("SELECT username from users where username=$1", user)
	if err != nil {
		return ctx, ErrUnauthorized
	}
	for row.Next() {
		if err := row.Scan(&value); err != nil {
			return ctx, ErrUnauthorized
		}
	}
	if value != "" {
		return ctx, ErrUnauthorized
	}

	return ctx, nil
}

func (s *Service) Signup(ctx context.Context, p *signup.SignupPayload) (*signup.SignupResult, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	hashedPassword := auth.ShaHashing(p.Password)
	userID := uuid.New().String()
	_, err = dbPool.Query("INSERT INTO Users Values ($1, $2, $3, $4, $5, $6, $7)", userID, p.Username, p.User.FirstName, p.User.LastName, p.User.Phone, p.User.Email, hashedPassword)
	if err != nil {
		return nil, err
	}
	token, err := auth.MakeToken(p.Username, userID)
	if err != nil {
		return nil, err
	}
	resp := signup.SignupResponse{JWT: &token, UserID: &userID}
	return &signup.SignupResult{User: &resp}, nil
}
