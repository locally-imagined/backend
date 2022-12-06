package auth

import (
	"fmt"
	"testing"
)

func TestMakeToken(t *testing.T) {
	type args struct {
		username string
		userID   string
	}
	want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6Imhpb29vIiwiVXNlcklEIjoiSGkiLCJpYXQiOjE2Njg1NzA3NjIsImp0aSI6IjI2NjliZDg5LTc3N2MtNDEyOC04NmZmLWFjNzRhNjBjZGI3ZSJ9.BZytXMptlcArcw_4HSBdxnpzN4J1XLqaSRlqh2A5WaY"
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{

			name:    "fail",
			args:    args{username: "hiooo", userID: "Hi"},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeToken(tt.args.username, tt.args.userID)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MakeToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
