package upload

import (
	"backend/auth"
	"backend/gen/upload"
	"context"

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
	// return context.WithValue(ctx, ctxValueClaims, auth), nil
	return ctx, nil
}

func (s *Service) UploadPhoto(ctx context.Context, p *upload.UploadPhotoPayload) (*upload.UploadPhotoResult, error) {
	t := true
	star := "*"
	return &upload.UploadPhotoResult{Success: &t, AccessControlAllowOrigin: &star}, nil
}
