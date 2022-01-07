package youtube

import (
	"strconv"
	"strings"
)

type FormatList []Format

// Type returns a new FormatList filtered by mime type of video
func (list FormatList) Type(t string) (result FormatList) {
	for i := range list {
		if strings.Contains(list[i].MimeType, t) {
			result = append(result, list[i])
		}
	}
	return result
}

// Quality returns a new FormatList filtered by quality, quality label or itag,
func (list FormatList) Quality(quality string) (result FormatList) {
	for _, f := range list {
		itag, _ := strconv.Atoi(quality)
		if itag == f.ItagNo || strings.Contains(f.Quality, quality) || strings.Contains(f.QualityLabel, quality) {
			result = append(result, f)
		}
	}
	return result
}
