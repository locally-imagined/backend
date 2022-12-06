package postings

import (
	"backend/gen/postings"
	"context"
	"testing"
)

func TestCreatePost(t *testing.T) {
	var (
		title        = "New Post"
		description  = "This describes the post"
		price        = "999"
		medium       = "Painting"
		imageIDs     = []string{"imageID1", "imageID2"}
		postID       = "postID"
		uploadDate   = "2022-10-09"
		sold         = false
		deliveryType = "Local Delivery"
		userID       = "userID"
		content      = []string{"MyBinaryImageData1, MyBinaryImageData2"}
		token        = "MyJWT"

		resp = postings.PostResponse{
			Title:        title,
			Description:  description,
			Price:        price,
			ImageIDs:     imageIDs,
			PostID:       postID,
			Medium:       medium,
			UploadDate:   uploadDate,
			Sold:         sold,
			Deliverytype: deliveryType,
			UserID:       userID,
		}
		res = postings.CreatePostResult{
			Posted: &resp,
		}
	)

	cases := []struct {
		Name        string
		Payload     *postings.CreatePostPayload
		Expected    *postings.CreatePostResult
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Payload:     MakeCreatePostPayload(t, token, title, description, price, medium, deliveryType, content),
			Expected:    &res,
			ExpectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mock := NewMock(t)
			mock.SetCreatePostFunc(func(ctx context.Context, p *postings.CreatePostPayload) (*postings.CreatePostResult, error) {
				return &res, nil
			})
			svc := NewService(mock)

			_, err := svc.CreatePost(context.Background(), c.Payload)
			if err != nil {
				t.Errorf("should not be erroring, this is a simple test")
				return
			}
			if res.Posted.Title != "New Post" {
				t.Errorf("got %s results, expected %s", res.Posted.Title, title)
				return
			}
			if res.Posted.Description != "This describes the post" {
				t.Errorf("got %s results, expected %s", res.Posted.Description, description)
				return
			}
			if res.Posted.Description != "This describes the post" {
				t.Errorf("got %s results, expected %s", res.Posted.Description, description)
				return
			}
			if res.Posted.ImageIDs[1] != "imageID2" {
				t.Errorf("got %s results, expected %s", res.Posted.ImageIDs[1], imageIDs[1])
				return
			}

		})
	}
}

func MakeCreatePostPayload(t *testing.T, token, title, description, price, medium, deliverytype string, content []string) *postings.CreatePostPayload {
	t.Helper()
	post := postings.Post{
		Title:        title,
		Description:  description,
		Price:        price,
		Medium:       medium,
		Deliverytype: deliverytype,
		Content:      content,
	}
	return &postings.CreatePostPayload{
		Post:  &post,
		Token: token,
	}
}

func TestGetPostPage(t *testing.T) {

	cases := []struct {
		Name        string
		Payload     *postings.GetPostPagePayload
		Expected    *postings.GetPostPageResult
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Payload:     MakeGetPostPagePayload(t, 0),
			Expected:    &res,
			ExpectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mock := NewMock(t)
			mock.SetGetPostPageFunc(func(ctx context.Context, p *postings.GetPostPagePayload) (*postings.GetPostPageeResult, error) {
				return &res, nil
			})
			svc := NewService(mock)

			_, err := svc.GetPostPage(context.Background(), c.Payload)
			if err != nil {
				t.Errorf("should not be erroring, this is a simple test")
				return
			}
			if res.Posted.Title != "New Post" {
				t.Errorf("got %s results, expected %s", res.Posted.Title, title)
				return
			}
			if res.Posted.Description != "This describes the post" {
				t.Errorf("got %s results, expected %s", res.Posted.Description, description)
				return
			}
			if res.Posted.Description != "This describes the post" {
				t.Errorf("got %s results, expected %s", res.Posted.Description, description)
				return
			}
			if res.Posted.ImageIDs[1] != "imageID2" {
				t.Errorf("got %s results, expected %s", res.Posted.ImageIDs[1], imageIDs[1])
				return
			}

		})
	}
}

func MakeGetPostPagePayload(t *testing.T, page int) *postings.GetPostPagePayload {
	t.Helper()
	return &postings.GetPostPagePayload{
		Page: page,
	}
}
