package login

import (
	"backend/gen/auth"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/mysql/mysql"
)

type Service struct{}

func (s *Service) Login(ctx context.Context, p *auth.LoginPayload) (string, error) {
	cleanup, err := mysql.RegisterDriver("cloudsql-mysql", cloudsqlconn.WithCredentialsFile("key.json"))
	return "hi", err
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbInstance := os.Getenv("INSTANCE")
	dbName := os.Getenv("DBNAME")
	conn := dbUser + ":" + dbPass + "@cloudsql-mysql(" + dbInstance + ")/" + dbName
	db, err := sql.Open(
		"cloudsql-mysql",
		conn,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT user, pass FROM test_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return "hi", err
}
