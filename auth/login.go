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
)

type Service struct{}

func (s *Service) Login(ctx context.Context, p *auth.LoginPayload) (string, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", fmt.Errorf("sql.Open: %v", err)
	}
	defer dbPool.Close()
	var password string
	hashedPassword := shaHashing(*p.Password)
	// Query for a value based on a single row.
	row, err := dbPool.Query("SELECT password from test_users where username=$1", *p.Username)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("account not found")
	}
	// return "", fmt.Errorf("%w", row)
	for row.Next() {
		if err := row.Scan(&password); err != nil {
			log.Fatal(err)
		}
	}

	if hashedPassword != password {
		return "BADPASSWORD", nil
	}
	token, err := MakeToken(*p.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Signup(ctx context.Context, p *auth.SignupPayload) (string, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", fmt.Errorf("sql.Open: %v", err)
	}
	defer dbPool.Close()
	hashedPassword := shaHashing(*p.Password)
	var value string = ""
	// Query for a value based on a single row.
	row, err := dbPool.Query("SELECT username from test_table where username=$1", *p.Username)
	if err != nil {
		return "", err
	}
	for row.Next() {
		if err := row.Scan(&value); err != nil {
			log.Fatal(err)
		}
	}
	if value != "" {
		return "", fmt.Errorf("account already exists")
	}

	_, err = dbPool.Query("INSERT INTO test_table (username, password) Values ('" + *p.Username + "', '" + hashedPassword + "')")
	if err != nil {
		return "", fmt.Errorf("account creation failed")
	}
	token, err := MakeToken(*p.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func shaHashing(input string) string {
	plainText := []byte(input)
	sha256Hash := sha256.Sum256(plainText)
	return hex.EncodeToString(sha256Hash[:])
}
