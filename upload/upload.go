package upload

import (
	"backend/auth"
	"backend/gen/upload"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"
	"goa.design/goa/v3/security"
)

type Service struct{}

var (
	// ErrUnauthorized is the error returned by Login when the request credentials
	// are invalid.
	ErrUnauthorized error = upload.Unauthorized("invalid jwt")
)

func (s *Service) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	tok := auth.DecodeToken(token)
	if tok == nil {
		return ctx, ErrUnauthorized
	}
	// 3. add authInfo to context

	claims := tok.Claims.(jwt.MapClaims)
	ctx = context.WithValue(ctx, "Name", claims["Name"])
	var tm time.Time
	switch iat := claims["iat"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	}
	ctx = context.WithValue(ctx, "exptime", tm)
	return ctx, nil
}

func (s *Service) UploadPhoto(ctx context.Context, p *upload.UploadPhotoPayload) (*upload.UploadPhotoResult, error) {
	var name string = ctx.Value("exptime").(string)
	star := "*"
	return &upload.UploadPhotoResult{Success: &name, AccessControlAllowOrigin: &star}, nil
}
