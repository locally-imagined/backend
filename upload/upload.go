package upload

import (
	"backend/gen/upload"
	"context"

	_ "github.com/lib/pq"
)

type Service struct{}

func (s *Service) UploadPhoto(ctx context.Context, p *upload.UploadPhotoPayload) (*upload.UploadPhotoResult, error) {
	t := true
	star := "*"
	return &upload.UploadPhotoResult{Success: &t, AccessControlAllowOrigin: &star}, nil
}
