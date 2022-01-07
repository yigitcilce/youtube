package youtube

import (
	"testing"
)

const (
	dwlURL    string = "https://www.youtube.com/watch?v=rFejpH_tAHM"
	streamURL string = "https://www.youtube.com/watch?v=5qap5aO4i9A"
	errURL    string = "https://www.youtube.com/watch?v=I8oGsuQ"
)

func TestYoutube_findVideoID(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid url",
			args: args{
				dwlURL,
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "valid id",
			args: args{
				"rFejpH_tAHM",
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "invalid character in id",
			args: args{
				"<M1!3",
			},
			wantErr:     true,
			expectedErr: ErrInvalidCharactersInVideoID,
		},
		{
			name: "video id is less than 10 characters",
			args: args{
				"afasda",
			},
			wantErr:     true,
			expectedErr: ErrVideoIDMinLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ExtractVideoID(tt.args.url); (err != nil) != tt.wantErr || err != tt.expectedErr {
				t.Errorf("extractVideoID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
