package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type locallyImaginedClaims struct {
	Name  string
	Email string
	jwt.RegisteredClaims
}

// change token sign from 'test'
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
	signedJWT, err := token.SignedString([]byte("test"))
	if err != nil {
		return "", fmt.Errorf("unable to sign token with secret: %v", err)
	}
	return signedJWT, nil
}

func DecodeToken(tokenString string) *jwt.Token {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("test"), nil
	})
	if err != nil {
		return nil
	}
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
	}
	return token
}

func ShaHashing(input string) string {
	plainText := []byte(input)
	sha256Hash := sha256.Sum256(plainText)
	return hex.EncodeToString(sha256Hash[:])
}
