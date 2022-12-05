package users

import (
	"backend/auth"
	"backend/gen/users"
	"context"

	_ "github.com/lib/pq"
	"goa.design/goa/v3/security"
)

type Service struct {
	usersClient Client
}

func NewService(usersClient Client) *Service {
	return &Service{
		usersClient: usersClient,
	}
}

func (s *Service) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return auth.JWTAuth(ctx, token, scheme)
}

func (s *Service) UpdateBio(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error) {
	res, err := s.usersClient.UpdateBio(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetContactInfo(ctx context.Context, p *users.GetContactInfoPayload) (*users.GetContactInfoResult, error) {
	res, err := s.usersClient.GetContactInfo(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) UpdateProfilePicture(ctx context.Context, p *users.UpdateProfilePicturePayload) (*users.UpdateProfilePictureResult, error) {
	res, err := s.usersClient.UpdateProfilePicture(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}
