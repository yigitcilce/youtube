package youtube

import "fmt"

type constError string
type ErrPlayabiltyStatus struct {
	Status string
	Reason string
}

const (
	ErrCipherNotFound             = constError("cipher not found")
	ErrInvalidCharactersInVideoID = constError("invalid characters in video id")
	ErrSignatureTimestampNotFound = constError("signature timestamp not found")
	ErrVideoPrivate               = constError("user restricted access to this video")
	ErrVideoIDMinLength           = constError("the video id must be at least 10 characters long")
)

func (e constError) Error() string {
	return string(e)
}

func (err ErrPlayabiltyStatus) Error() string {
	return fmt.Sprintf("cannot playback and download, status: %s, reason: %s", err.Status, err.Reason)
}

// ErrUnexpectedHTTPStatusCode is returned on unexpected HTTP status codes
type ErrUnexpectedHTTPStatusCode int

func (err ErrUnexpectedHTTPStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code: %d", err)
}
