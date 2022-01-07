package youtube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type contentPlaybackContext struct {
	SignatureTimestamp string `json:"signatureTimestamp"`
}
type playbackContext struct {
	ContentPlaybackContext contentPlaybackContext `json:"contentPlaybackContext"`
}

type inntertubeContext struct {
	Client innertubeClient `json:"client"`
}
type innertubeRequest struct {
	VideoID         string            `json:"videoId"`
	Context         inntertubeContext `json:"context"`
	PlaybackContext playbackContext   `json:"playbackContext"`
}

type innertubeClient struct {
	HL            string `json:"hl"`
	GL            string `json:"gl"`
	ClientName    string `json:"clientName"`
	ClientVersion string `json:"clientVersion"`
}

type ClientType string

// Client offers methods to download video metadata and video streams.
type Client struct {
	// HTTPClient can be used to set a custom HTTP client.
	// If not set, http.DefaultClient will be used
	HTTPClient *http.Client

	// playerCache caches the JavaScript code of a player response
	playerCache playerCache
}

// GetVideo fetches video metadata
func (c *Client) GetVideo(url string) (*Video, error) {
	return c.GetVideoContext(context.Background(), url)
}

// GetVideoContext fetches video metadata with a context
func (c *Client) GetVideoContext(ctx context.Context, url string) (*Video, error) {
	id, err := ExtractVideoID(url)
	if err != nil {
		return nil, fmt.Errorf("extractVideoID failed: %w", err)
	}
	return c.videoFromID(ctx, id)
}

func (c *Client) videoFromID(ctx context.Context, id string) (*Video, error) {
	body, err := c.videoDataByInnertube(ctx, id, Web)
	if err != nil {
		return nil, err
	}

	v := &Video{
		ID: id,
	}

	err = v.parseVideoInfo(body)
	if err == nil {
		return v, nil
	}

	// error while parsing
	return v, err
}

const Web ClientType = "WEB"

func (c *Client) videoDataByInnertube(ctx context.Context, id string, clientType ClientType) ([]byte, error) {
	config, err := c.getPlayerConfig(ctx, id)
	if err != nil {
		return nil, err
	}

	// fetch sts first
	sts, err := config.getSignatureTimestamp()
	if err != nil {
		return nil, err
	}

	data, keyToken := prepareInnertubeData(id, sts, clientType)
	reqData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("https://www.youtube.com/youtubei/v1/player?key=%s", keyToken)

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return io.ReadAll(resp.Body)
}

var innertubeClientInfo = map[ClientType]map[string]string{
	Web: {
		"version": "2.20210617.01.00",
		"key":     "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8",
	},
}

func prepareInnertubeData(videoID string, sts string, clientType ClientType) (innertubeRequest, string) {
	cInfo, ok := innertubeClientInfo[clientType]
	if !ok {
		// if provided clientType not exist - use Web as fallback option
		clientType = Web
		cInfo = innertubeClientInfo[clientType]
	}

	return innertubeRequest{
		VideoID: videoID,
		Context: inntertubeContext{
			Client: innertubeClient{
				HL:            "en",
				GL:            "US",
				ClientName:    string(clientType),
				ClientVersion: cInfo["version"],
			},
		},
		PlaybackContext: playbackContext{
			ContentPlaybackContext: contentPlaybackContext{
				SignatureTimestamp: sts,
			},
		},
	}, cInfo["key"]
}

// GetStream returns the stream and the total size for a specific format
func (c *Client) GetStream(video *Video, format *Format) (io.ReadCloser, int64, error) {
	return c.GetStreamContext(context.Background(), video, format)
}

// GetStreamContext returns the stream and the total size for a specific format with a context
func (c *Client) GetStreamContext(ctx context.Context, video *Video, format *Format) (io.ReadCloser, int64, error) {
	url, err := c.GetStreamURL(video, format)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	r, w := io.Pipe()

	// go magic starts here
	go c.download(req, w, format)

	return r, format.ContentLength, nil
}

// download gets the http response body and writes into memory
func (c *Client) download(req *http.Request, w *io.PipeWriter, format *Format) {
	const chunkSize int64 = 10000000

	// Get http body content by pieces till nothing is left
	loadChunk := func(pos int64) (int64, error) {
		req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", pos, pos+chunkSize-1))

		resp, err := c.httpDo(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusPartialContent {
			return 0, ErrUnexpectedHTTPStatusCode(resp.StatusCode)
		}

		return io.Copy(w, resp.Body)
	}
	defer w.Close()

	if format.ContentLength == 0 {
		resp, err := c.httpDo(req)
		if err != nil {
			w.CloseWithError(err)
			return
		}

		defer resp.Body.Close()

		io.Copy(w, resp.Body)
		return
	}

	/// Downloading in multiple chunks is much faster
	for pos := int64(0); pos < format.ContentLength; {
		written, err := loadChunk(pos)
		if err != nil {
			w.CloseWithError(err)
			return
		}

		pos += written
	}
}

// GetStreamURL returns the url for a specific format
func (c *Client) GetStreamURL(video *Video, format *Format) (string, error) {
	return c.GetStreamURLContext(context.Background(), video, format)
}

// GetStreamURLContext returns the url for a specific format with a context
func (c *Client) GetStreamURLContext(ctx context.Context, video *Video, format *Format) (string, error) {
	cipher := format.Cipher
	if cipher == "" {
		return "", ErrCipherNotFound
	}

	// don't know whats going in decipher, goodluck reading
	uri, err := c.decipherURL(ctx, video.ID, cipher)
	if err != nil {
		return "", err
	}

	return uri, err
}

// httpDo sends an HTTP request and returns an HTTP response.
func (c *Client) httpDo(req *http.Request) (*http.Response, error) {
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)

	if res != nil {
		log.Println(res.Status)
	}

	return res, err
}

// httpGet does a HTTP GET request, checks the response to be a 200 OK and returns it
func (c *Client) httpGet(ctx context.Context, url string) (*http.Response, error) {
	// Prepare GET with given context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Sents the HTTP request, waits the response
	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}

	// If there is a problem with http response, return it
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, ErrUnexpectedHTTPStatusCode(resp.StatusCode)
	}

	return resp, nil
}

// httpGetBodyBytes reads the whole HTTP body and returns it
func (c *Client) httpGetBodyBytes(ctx context.Context, url string) ([]byte, error) {
	resp, err := c.httpGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
