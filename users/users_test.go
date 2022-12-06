package users

import (
	"backend/gen/users"
	"context"
	"testing"
)

func TestUpdateBio(t *testing.T) {
	var (
		firstname = "First"
		lastname  = "last"
		phone     = "999"
		email     = "a@test.com"
		bio       = "My Bio"
		ppid      = "ppid"
		token     = "MyJWT"

		resp = users.User{
			FirstName: firstname,
			LastName:  lastname,
			Phone:     phone,
			Email:     email,
			Bio:       bio,
			ProfpicID: ppid,
		}
		res = users.UpdateBioResult{
			User: &resp,
		}
	)

	cases := []struct {
		Name        string
		Payload     *users.UpdateBioPayload
		Expected    *users.UpdateBioResult
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Payload:     MakeUpdateBioPayload(t, token, bio),
			Expected:    &res,
			ExpectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mock := NewMock(t)
			mock.SetUpdateBioFunc(func(ctx context.Context, p *users.UpdateBioPayload) (*users.UpdateBioResult, error) {
				return &res, nil
			})
			svc := NewService(mock)

			_, err := svc.UpdateBio(context.Background(), c.Payload)
			if err != nil {
				t.Errorf("should not be erroring, this is a simple test")
				return
			}
			if res.User.FirstName != "First" {
				t.Errorf("got %s results, expected %s", res.User.FirstName, firstname)
				return
			}
			if res.User.Email != "a@test.com" {
				t.Errorf("got %s results, expected %s", res.User.Email, email)
				return
			}
		})
	}
}

func MakeUpdateBioPayload(t *testing.T, token, bio string) *users.UpdateBioPayload {
	t.Helper()
	return &users.UpdateBioPayload{
		Bio:   &users.Bio{Bio: &bio},
		Token: token,
	}
}

func TestUpdateProfilePicture(t *testing.T) {
	var (
		ppid    = "ppid"
		token   = "MyJWT"
		content = "binaryimagecontent"

		res = users.UpdateProfilePictureResult{
			ImageID: &users.ProfilePhoto{ImageID: &ppid},
		}
	)

	cases := []struct {
		Name        string
		Payload     *users.UpdateProfilePicturePayload
		Expected    *users.UpdateProfilePictureResult
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Payload:     MakeUpdateProfilePicturePayload(t, token, content),
			Expected:    &res,
			ExpectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mock := NewMock(t)
			mock.SetUpdateProfilePictureFunc(func(ctx context.Context, p *users.UpdateProfilePicturePayload) (*users.UpdateProfilePictureResult, error) {
				return &res, nil
			})
			svc := NewService(mock)

			_, err := svc.UpdateProfilePicture(context.Background(), c.Payload)
			if err != nil {
				t.Errorf("should not be erroring, this is a simple test")
				return
			}
			if *res.ImageID.ImageID != "ppid" {
				t.Errorf("got %s results, expected %s", *res.ImageID.ImageID, ppid)
				return
			}
		})
	}
}

func MakeUpdateProfilePicturePayload(t *testing.T, token, content string) *users.UpdateProfilePicturePayload {
	t.Helper()
	return &users.UpdateProfilePicturePayload{
		Content: &users.Content{Content: &content},
		Token:   token,
	}
}
