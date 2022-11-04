package upload

import (
	"backend/gen/upload"
	"context"
	"log"

	_ "github.com/lib/pq"
)

type Service struct{}

func (s *Service) UploadPhoto(ctx context.Context, p *upload.UploadPhotoPayload) (*upload.UploadPhotoResult, error) {
	t := true
	star := "*"
	log.Fatalf(*p.Authorization)
	return &upload.UploadPhotoResult{Success: &t, AccessControlAllowOrigin: &star}, nil
}
