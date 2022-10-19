package main

import "testing"

func TestMakeToken(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeToken(tt.args.email)
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
