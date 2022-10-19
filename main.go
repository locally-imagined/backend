package main

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func MakeToken(ctx context.Context, email string) (string, error) {
	uuid := uuid.New()
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       uuid.String()}

	payload := &zendeskClaims{
		Name:             email,
		Email:            email,
		RegisteredClaims: claims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedJWT, err := token.SignedString([]byte(c.sharedSecret))
	if err != nil {
		return "", fmt.Errorf("unable to sign token with zendesk secret: %v", err)
	}
	return signedJWT, nil
}

func main() {
	return

}
