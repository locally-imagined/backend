package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"backend/auth/helpers.go"

	_ "github.com/lib/pq"
)

type Service struct{}

type post struct {
	postID int
	postTitle string
	postDesc string
	userName string
	email string
}
  

func (s *Service) getAllPosts() ([]post, error) {
	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", fmt.Errorf("sql.Open: %v", err)
	}

	rows, err := dbPool.Query("SELECT * FROM posts LEFT JOIN images ON posts.postID = images.postID LEFT JOIN users ON posts.ownerID = users.userId")

	defer dbPool.Close()
	
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no posts available")
	}

	var allPosts []post

	// return "", fmt.Errorf("%w", row)
	for rows.Next() {
		var row post
		if err := rows.Scan(&row.postID, &row.postTitle, &row.postDesc, &row.userName, &row.email); err != nil {
			log.Fatal(err)
		}
		allPosts = append(allPosts, row)
	}
	
	return allPosts, nil
}

func (s *Service) getUserPosts(string jwt) ([]post, error) {
	token := DecodeToken(jwt)
	claims := token.Claims.(*locallyImaginedClaims) //NOTE: acting like userID is one of the claims


	dbPool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", fmt.Errorf("sql.Open: %v", err)
	}

	rows, err := dbPool.Query("SELECT * FROM posts WHERE ownerID = $1 LEFT JOIN images ON posts.postID = images.postID LEFT JOIN users ON posts.ownerID = users.userId", claims.userID)

	defer dbPool.Close()
	
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no posts available")
	}

	var allPosts []post

	// return "", fmt.Errorf("%w", row)
	for rows.Next() {
		var row post
		if err := rows.Scan(&row.postID, &row.postTitle, &row.postDesc, &row.userName, &row.email); err != nil {
			log.Fatal(err)
		}
		allPosts = append(allPosts, row)
	}
	
	return allPosts, nil
}