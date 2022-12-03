package users

import (
	"backend/gen/users"
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Service struct{}

func (s *Service) UpdateBio(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()
	_, err = dbPool.Query("UPDATE users SET bio = $1 WHERE userID = $2", p.Bio, ctx.Value("UserID").(string))
	if err != nil {
		return nil, err
	}

	rows, err := dbPool.Query("SELECT firstname, lastname, phone, email FROM users WHERE userID = $1", ctx.Value("UserID").(string))
	var resp users.UpdateBioResult
	for rows.Next() {
		var row users.User
		if err := rows.Scan(&row.FirstName, &row.LastName, &row.Phone, &row.Email); err != nil {
			log.Fatal(err)
			return nil, err
		}
		resp = users.UpdateBioResult{UpdatedUser: &row}
		//res = append(res, &users.User{FirstName: row.firstname, LastName: row.lastname, Phone: row.phone, Email: row.email})
	}
	return &resp, err
}

func (s *Service) GetContactInfo(ctx context.Context, p *users.GetContactInfoPayload) (*users.GetContactInfoResult, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer dbPool.Close()

	var rows *sql.Rows
	rows, err = dbPool.Query("SELECT firstname, lastname, phone, email WHERE userid = $1", p.Userid)
	if err != nil {
		return nil, err
	}

	var resp users.GetContactInfoResult
	for rows.Next() {
		var row users.User
		if err := rows.Scan(&row.FirstName, &row.LastName, &row.Phone, &row.Email); err != nil {
			log.Fatal(err)
			return nil, err
		}
		resp = users.GetContactInfoResult{ContactInfo: &row}
		//res = append(res, &users.User{FirstName: row.firstname, LastName: row.lastname, Phone: row.phone, Email: row.email})
	}
	return &resp, err
}
