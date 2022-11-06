package auth

import (
	"backend/gen/auth"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"goa.design/goa/v3/security"
)

type Service struct{}

var (
	// ErrUnauthorized is the error returned by Login when the request credentials
	// are invalid.
	ErrUnauthorized error = auth.Unauthorized("invalid username and password combination")
)

func (s *Service) BasicAuth(ctx context.Context, user, pass string, scheme *security.BasicScheme) (context.Context, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return ctx, ErrUnauthorized
	}
	defer dbPool.Close()
	var password string
	hashedPassword := shaHashing(pass)
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

func (s *Service) Login(ctx context.Context, p *auth.LoginPayload) (*auth.LoginResult, error) {
	access := "*"
	// dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	return &auth.LoginResult{JWT: nil, AccessControlAllowOrigin: &access}, fmt.Errorf("sql.Open: %v", err)
	// }
	// defer dbPool.Close()
	// var password string
	// hashedPassword := shaHashing(*p.Password)
	// // Query for a value based on a single row.
	// row, err := dbPool.Query("SELECT password from test_users where username=$1", *p.Username)
	// if err == sql.ErrNoRows {
	// 	return &auth.LoginResult{JWT: nil, AccessControlAllowOrigin: &access}, fmt.Errorf("account not found")
	// }
	// // return "", fmt.Errorf("%w", row)
	// for row.Next() {
	// 	if err := row.Scan(&password); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// if hashedPassword != password {
	// 	return &auth.LoginResult{JWT: nil, AccessControlAllowOrigin: &access}, nil
	// }
	token, err := MakeToken(p.Username)
	if err != nil {
		return &auth.LoginResult{JWT: nil, AccessControlAllowOrigin: &access}, err
	}
	return &auth.LoginResult{JWT: &token, AccessControlAllowOrigin: &access}, nil
}

func (s *Service) Signup(ctx context.Context, p *auth.SignupPayload) (*auth.SignupResult, error) {
	access := "*"
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return &auth.SignupResult{JWT: nil, AccessControlAllowOrigin: &access}, fmt.Errorf("sql.Open: %v", err)
	}
	defer dbPool.Close()
	hashedPassword := shaHashing(*p.Password)
	var value string = ""
	// Query for a value based on a single row.
	row, err := dbPool.Query("SELECT username from test_users where username=$1", *p.Username)
	if err != nil {
		return &auth.SignupResult{JWT: nil, AccessControlAllowOrigin: &access}, err
	}
	for row.Next() {
		if err := row.Scan(&value); err != nil {
			log.Fatal(err)
		}
	}
	if value != "" {
		return &auth.SignupResult{JWT: nil, AccessControlAllowOrigin: &access}, fmt.Errorf("account already exists")
	}

	// double check this
	_, err = dbPool.Query("INSERT INTO test_users (username, password) Values ($1, $2)", *p.Username, hashedPassword)
	if err != nil {
		return &auth.SignupResult{JWT: nil, AccessControlAllowOrigin: &access}, fmt.Errorf("account creation failed")
	}
	token, err := MakeToken(*p.Username)
	if err != nil {
		return &auth.SignupResult{JWT: nil, AccessControlAllowOrigin: &access}, err
	}
	return &auth.SignupResult{JWT: &token, AccessControlAllowOrigin: &access}, nil
}

func shaHashing(input string) string {
	plainText := []byte(input)
	sha256Hash := sha256.Sum256(plainText)
	return hex.EncodeToString(sha256Hash[:])
}
