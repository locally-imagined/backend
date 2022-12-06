package users

import (
	"backend/gen/users"
	"context"
	"testing"

	"goa.design/clue/mock"
)

type (
	UpdateBioFunc            func(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error)
	GetUserInfoFunc          func(ctx context.Context, p *users.GetUserInfoPayload) (*users.GetUserInfoResult, error)
	UpdateProfilePictureFunc func(ctx context.Context, p *users.UpdateProfilePicturePayload) (*users.UpdateProfilePictureResult, error)
	Mock                     struct {
		m *mock.Mock
		t *testing.T
	}
)

var _ Client = &Mock{}

func NewMock(t *testing.T) *Mock {
	return &Mock{mock.New(), t}
}

func (m *Mock) AddUpdateBioFunc(f UpdateBioFunc) {
	m.m.Add("UpdateBio", f)
}

func (m *Mock) AddGetUserInfoFunc(f GetUserInfoFunc) {
	m.m.Add("GetUserInfo", f)
}

func (m *Mock) AddUpdateProfilePictureFunc(f UpdateProfilePictureFunc) {
	m.m.Add("UpdateProfilePicture", f)
}

func (m *Mock) SetUpdateBioFunc(f UpdateBioFunc) {
	m.m.Set("UpdateBio", f)
}
func (m *Mock) SetGetUserInfoFunc(f GetUserInfoFunc) {
	m.m.Set("GetUserInfo", f)
}
func (m *Mock) SetUpdateProfilePictureFunc(f UpdateProfilePictureFunc) {
	m.m.Set("UpdateProfilePicture", f)
}

func (m *Mock) UpdateBio(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error) {
	if f := m.m.Next("UpdateBio"); f != nil {
		return f.(UpdateBioFunc)(ctx, p)
	}
	m.t.Error("unexpected call to UpdateBio")
	return nil, nil
}

func (m *Mock) GetUserInfo(ctx context.Context, p *users.GetUserInfoPayload) (*users.GetUserInfoResult, error) {
	if f := m.m.Next("GetUserInfo"); f != nil {
		return f.(GetUserInfoFunc)(ctx, p)
	}
	m.t.Error("unexpected call to GetUserInfo")
	return nil, nil
}

func (m *Mock) UpdateProfilePicture(ctx context.Context, p *users.UpdateProfilePicturePayload) (*users.UpdateProfilePictureResult, error) {
	if f := m.m.Next("UpdateProfilePicture"); f != nil {
		return f.(UpdateProfilePictureFunc)(ctx, p)
	}
	m.t.Error("unexpected call to UpdateProfilePicture")
	return nil, nil
}

func (m *Mock) HasMore() bool {
	return m.m.HasMore()
}
