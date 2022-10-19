package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type locallyImaginedClaims struct {
	Name  string
	Email string
	jwt.RegisteredClaims
}

func MakeToken(email string) (string, error) {
	uuid := uuid.New()
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       uuid.String()}

	payload := &locallyImaginedClaims{
		Name:             email,
		Email:            email,
		RegisteredClaims: claims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	log.Println(token)
	signedJWT, err := token.SignedString([]byte("test"))
	if err != nil {
		return "", fmt.Errorf("unable to sign token with zendesk secret: %v", err)
	}
	return signedJWT, nil
}

func DecodeToken(token string) string {
	return "test"
}

func main() {
	return

}
