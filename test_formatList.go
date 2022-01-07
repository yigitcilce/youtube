package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatList_Type(t *testing.T) {
	list := []Format{{
		MimeType: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"",
	},
	}
	type args struct {
		mimeType string
	}
	tests := []struct {
		name string
		list FormatList
		args args
		want FormatList
	}{
		{
			name: "find video",
			list: list,
			args: args{
				mimeType: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"",
			},
			want: []Format{{
				MimeType: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := tt.list.Type("video")
			assert.Equal(t, format, tt.want)
		})
	}
}
