package main

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/yigitcilce/youtube"
)

var downloader *Downloader

func getDownloader() *Downloader {
	// Once downloader is configured, use it
	if downloader != nil {
		return downloader
	}

	// Initialize Downloader
	downloader = &Downloader{}

	// Connection rules for downloader
	httpTransport := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	// Assign http rules
	downloader.HTTPClient = &http.Client{Transport: httpTransport}

	return downloader
}

// getVideoWithFormat gets video and its format for downloading process
func getVideoWithFormat(id string) (*youtube.Video, *youtube.Format, error) {
	yt := getDownloader()
	video, err := yt.GetVideo(id)
	if err != nil {
		return nil, nil, err
	}
	formats := video.Formats

	if len(formats) == 0 {
		return nil, nil, errors.New("no formats found")
	}
	return video, &formats[0], nil
}
