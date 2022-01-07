package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// For extracting VideoID
var videoRegexpList = []*regexp.Regexp{
	regexp.MustCompile(`(?:v|embed|shorts|watch\?v)(?:=|/)([^"&?/=%]{11})`),
	regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
	regexp.MustCompile(`([^"&?/=%]{11})`),
}

type Video struct {
	ID          string
	Title       string
	Description string
	Duration    time.Duration
	PublishDate time.Time
	Formats     FormatList
}

// parseVideoInfo parses video information from http response body
func (v *Video) parseVideoInfo(body []byte) error {
	var prData playerResponseData

	err := json.Unmarshal(body, &prData)
	if err != nil {
		return fmt.Errorf("unable to parse player response JSON: %w", err)
	}

	return v.extractDataFromPlayerResponse(prData)
}

func (v *Video) extractDataFromPlayerResponse(prData playerResponseData) error {
	// Get title for file creation
	v.Title = prData.VideoDetails.Title

	// Assign Streams for download process
	v.Formats = append(prData.StreamingData.Formats, prData.StreamingData.AdaptiveFormats...)
	if len(v.Formats) == 0 {
		return errors.New("no formats found in the server's answer")
	}
	return nil
}

// ExtractVideoID extracts the videoID from the given string
func ExtractVideoID(videoID string) (string, error) {
	if strings.ContainsAny(videoID, "\"?&/<%=") {
		for _, re := range videoRegexpList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}

	if strings.ContainsAny(videoID, "?&/<%=") {
		return "", ErrInvalidCharactersInVideoID
	}
	if len(videoID) < 10 {
		return "", ErrVideoIDMinLength
	}

	return videoID, nil
}
