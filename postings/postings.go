package postings

// TODO
// add delivery type column to db and add support for it in createpost and editpost endpoints
// add mark as sold option to ui and change editpost endpoint to support it
// add artist bio column to users default empty and create endpoint to set user bio
// decide on filters: art type, price, sold?

// should be package postings
import (
	"backend/auth"
	"backend/gen/postings"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"

	_ "github.com/lib/pq"
	"goa.design/goa/v3/security"
)

type Service struct {
	postingsClient Client
}

func NewService(postingsClient Client) *Service {
	return &Service{
		postingsClient: postingsClient,
	}
}

func (s *Service) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	tok := auth.DecodeToken(token)
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

func (s *Service) CreatePost(ctx context.Context, p *postings.CreatePostPayload) (*postings.CreatePostResult, error) {
	res, err := s.postingsClient.CreatePost(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetPostPage(ctx context.Context, p *postings.GetPostPagePayload) (*postings.GetPostPageResult, error) {
	res, err := s.postingsClient.GetPostPage(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetArtistPostPage(ctx context.Context, p *postings.GetArtistPostPagePayload) (*postings.GetArtistPostPageResult, error) {
	res, err := s.postingsClient.GetArtistPostPage(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetPostPageFiltered(ctx context.Context, p *postings.GetPostPageFilteredPayload) (*postings.GetPostPageFilteredResult, error) {
	res, err := s.postingsClient.GetPostPageFiltered(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetImagesForPost(ctx context.Context, p *postings.GetImagesForPostPayload) (*postings.GetImagesForPostResult, error) {
	res, err := s.postingsClient.GetImagesForPost(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) DeletePost(ctx context.Context, p *postings.DeletePostPayload) error {
	err := s.postingsClient.DeletePost(ctx, p)
	return err
}

func (s *Service) EditPost(ctx context.Context, p *postings.EditPostPayload) (*postings.EditPostResult, error) {
	res, err := s.postingsClient.EditPost(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}
