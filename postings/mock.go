package postings

import (
	"backend/gen/postings"
	"context"
	"database/sql"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"goa.design/clue/mock"
)

type (
	CreatePostFunc          func(ctx context.Context, p *postings.CreatePostPayload) (*postings.CreatePostResult, error)
	GetPostPageFunc         func(ctx context.Context, p *postings.GetPostPagePayload) (*postings.GetPostPageResult, error)
	GetArtistPostPageFunc   func(ctx context.Context, p *postings.GetArtistPostPagePayload) (*postings.GetArtistPostPageResult, error)
	GetPostPageFilteredFunc func(ctx context.Context, p *postings.GetPostPageFilteredPayload) (*postings.GetPostPageFilteredResult, error)
	GetImagesForPostFunc    func(ctx context.Context, p *postings.GetImagesForPostPayload) (*postings.GetImagesForPostResult, error)
	DeletePostFunc          func(ctx context.Context, p *postings.DeletePostPayload) error
	EditPostFunc            func(ctx context.Context, p *postings.EditPostPayload) (*postings.EditPostResult, error)
	openDBFunc              func() (*sql.DB, error)
	getS3SessionFunc        func() (string, *s3.S3)
	Mock                    struct {
		m *mock.Mock
		t *testing.T
	}
)

var _ Client = &Mock{}

func NewMock(t *testing.T) *Mock {
	return &Mock{mock.New(), t}
}

func (m *Mock) AddCreatePostFunc(f CreatePostFunc) {
	m.m.Add("CreatePost", f)
}
func (m *Mock) AddEditPostFunc(f EditPostFunc) {
	m.m.Add("EditPost", f)
}
func (m *Mock) AddDeletePostFunc(f DeletePostFunc) {
	m.m.Add("DeletePost", f)
}
func (m *Mock) AddGetPostPageFunc(f GetPostPageFunc) {
	m.m.Add("GetPostPage", f)
}
func (m *Mock) AddGetArtistPostPageFunc(f GetArtistPostPageFunc) {
	m.m.Add("GetArtistPostPage", f)
}
func (m *Mock) AddGetPostPageFilteredFunc(f GetPostPageFilteredFunc) {
	m.m.Add("GetPostPageFiltered", f)
}
func (m *Mock) AddGetImagesForPostFunc(f GetImagesForPostFunc) {
	m.m.Add("GetImagesForPost", f)
}
func (m *Mock) AddopenDBFunc(f openDBFunc) {
	m.m.Add("openDB", f)
}
func (m *Mock) AddgetS3SessionFunc(f getS3SessionFunc) {
	m.m.Add("getS3Session", f)
}

func (m *Mock) SetCreatePostFunc(f CreatePostFunc) {
	m.m.Set("CreatePost", f)
}
func (m *Mock) SetEditPostFunc(f EditPostFunc) {
	m.m.Set("EditPost", f)
}
func (m *Mock) SetDeletePostFunc(f DeletePostFunc) {
	m.m.Set("DeletePost", f)
}
func (m *Mock) SetGetPostPageFunc(f GetPostPageFunc) {
	m.m.Set("GetPostPage", f)
}
func (m *Mock) SetGetArtistPostPageFunc(f GetArtistPostPageFunc) {
	m.m.Set("GetArtistPostPage", f)
}
func (m *Mock) SetGetPostPageFilteredFunc(f GetPostPageFilteredFunc) {
	m.m.Set("GetPostPageFiltered", f)
}
func (m *Mock) SetGetImagesForPostFunc(f GetImagesForPostFunc) {
	m.m.Set("GetImagesForPost", f)
}
func (m *Mock) SetopenDBFunc(f openDBFunc) {
	m.m.Set("openDB", f)
}
func (m *Mock) SetgetS3SessionFunc(f getS3SessionFunc) {
	m.m.Set("getS3Session", f)
}

func (m *Mock) CreatePost(ctx context.Context, p *postings.CreatePostPayload) (*postings.CreatePostResult, error) {
	if f := m.m.Next("CreatePost"); f != nil {
		return f.(CreatePostFunc)(ctx, p)
	}
	m.t.Error("unexpected call to CreatePost")
	return nil, nil
}

func (m *Mock) EditPost(ctx context.Context, p *postings.EditPostPayload) (*postings.EditPostResult, error) {
	if f := m.m.Next("CreatePost"); f != nil {
		return f.(EditPostFunc)(ctx, p)
	}
	m.t.Error("unexpected call to EditPost")
	return nil, nil
}

func (m *Mock) DeletePost(ctx context.Context, p *postings.DeletePostPayload) error {
	if f := m.m.Next("DeletePost"); f != nil {
		return f.(DeletePostFunc)(ctx, p)
	}
	m.t.Error("unexpected call to DeletePost")
	return nil
}

func (m *Mock) GetPostPage(ctx context.Context, p *postings.GetPostPagePayload) (*postings.GetPostPageResult, error) {
	if f := m.m.Next("GetPostPage"); f != nil {
		return f.(GetPostPageFunc)(ctx, p)
	}
	m.t.Error("unexpected call to GetPostPagePost")
	return nil, nil
}

func (m *Mock) GetArtistPostPage(ctx context.Context, p *postings.GetArtistPostPagePayload) (*postings.GetArtistPostPageResult, error) {
	if f := m.m.Next("GetArtistPostPage"); f != nil {
		return f.(GetArtistPostPageFunc)(ctx, p)
	}
	m.t.Error("unexpected call to GetArtistPostPagePost")
	return nil, nil
}

func (m *Mock) GetPostPageFiltered(ctx context.Context, p *postings.GetPostPageFilteredPayload) (*postings.GetPostPageFilteredResult, error) {
	if f := m.m.Next("GetPostPageFiltered"); f != nil {
		return f.(GetPostPageFilteredFunc)(ctx, p)
	}
	m.t.Error("unexpected call to GetPostPageFilteredPost")
	return nil, nil
}

func (m *Mock) GetImagesForPost(ctx context.Context, p *postings.GetImagesForPostPayload) (*postings.GetImagesForPostResult, error) {
	if f := m.m.Next("GetImagesForPost"); f != nil {
		return f.(GetImagesForPostFunc)(ctx, p)
	}
	m.t.Error("unexpected call to GetImagesForPost")
	return nil, nil
}

func (m *Mock) openDB() (*sql.DB, error) {
	if f := m.m.Next("openDB"); f != nil {
		return f.(openDBFunc)()
	}
	m.t.Error("unexpected call to openDB")
	return nil, nil
}

func (m *Mock) getS3Session() (string, *s3.S3) {
	if f := m.m.Next("getS3Session"); f != nil {
		return f.(getS3SessionFunc)()
	}
	m.t.Error("unexpected call to getS3Session")
	return "", nil
}

func (m *Mock) HasMore() bool {
	return m.m.HasMore()
}
