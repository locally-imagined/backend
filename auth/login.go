package auth

import (
	"backend/gen/auth"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/go-sql-driver/mysql"
)

type Service struct{}

func (s *Service) Login(ctx context.Context, p *auth.LoginPayload) (string, error) {
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser                 = os.Getenv("DBUSER")   // e.g. 'my-db-user'
		dbPwd                  = os.Getenv("DBPASS")   // e.g. 'my-db-password'
		dbName                 = os.Getenv("DBNAME")   // e.g. 'my-database'
		instanceConnectionName = os.Getenv("INSTANCE") // e.g. 'project:region:instance'
	)

	d, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		return "", fmt.Errorf("cloudsqlconn.NewDialer: %v", err)
	}
	mysql.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, addr string) (net.Conn, error) {
			return d.Dial(ctx, instanceConnectionName)
		})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return "", fmt.Errorf("sql.Open: %v", err)
	}
	var password string
	hashedPassword := shaHashing(*p.Password)
	// Query for a value based on a single row.
	row, err := dbPool.Query("SELECT password from test_table where username='" + *p.Username + "'")
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
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser                 = os.Getenv("DBUSER")   // e.g. 'my-db-user'
		dbPwd                  = os.Getenv("DBPASS")   // e.g. 'my-db-password'
		dbName                 = os.Getenv("DBNAME")   // e.g. 'my-database'
		instanceConnectionName = os.Getenv("INSTANCE") // e.g. 'project:region:instance'
	)

	d, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		return "", fmt.Errorf("cloudsqlconn.NewDialer: %v", err)
	}
	mysql.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, addr string) (net.Conn, error) {
			return d.Dial(ctx, instanceConnectionName)
		})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return "", fmt.Errorf("sql.Open: %v", err)
	}
	hashedPassword := shaHashing(*p.Password)
	var value string = ""
	// Query for a value based on a single row.
	//row, err := dbPool.Query("INSERT INTO test_table (username, password) Values ('" + *p.Username + "', '" + hashedPassword + "'")
	row, err := dbPool.Query("SELECT username from test_table where username='" + *p.Username + "'")
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
