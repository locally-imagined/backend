package auth

import (
	"backend/gen/postings"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

type locallyImaginedClaims struct {
	Username string
	UserID   string
	jwt.RegisteredClaims
}

var ErrUnauthorized *goa.ServiceError = postings.MakeUnauthorized(fmt.Errorf("Unauthorized"))

func JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	tok := DecodeToken(token)
	if tok == nil {
		return ctx, ErrUnauthorized
	}
	// 3. add authInfo to context
	claims := tok.Claims.(jwt.MapClaims)
	ctx = context.WithValue(ctx, "Username", claims["Username"])
	ctx = context.WithValue(ctx, "UserID", claims["UserID"])
	var exp time.Time
	var now time.Time = time.Now()
	switch iat := claims["iat"].(type) {
	case float64:
		exp = time.Unix(int64(iat), 0)
	}
	if exp.Add(time.Hour * 2).Before(now) {
		return ctx, ErrUnauthorized
	}

	return ctx, nil
}

// change token sign from 'test'
func MakeToken(username, userID string) (string, error) {
	uuid := uuid.New()
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       uuid.String()}

	payload := &locallyImaginedClaims{
		Username:         username,
		UserID:           userID,
		RegisteredClaims: claims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedJWT, err := token.SignedString([]byte("test"))
	if err != nil {
		return "", ErrUnauthorized
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
