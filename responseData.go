package youtube

// responseData presents a part of youtubes video API
type playerResponseData struct {
	StreamingData struct {
		Formats         []Format `json:"formats"`
		AdaptiveFormats []Format `json:"adaptiveFormats"`
	} `json:"streamingData"`
	VideoDetails struct {
		Title string `json:"title"`
	} `json:"videoDetails"`
}

type Format struct {
	ItagNo   int    `json:"itag"`
	MimeType string `json:"mimeType"`
	Quality  string `json:"quality"`
	Cipher   string `json:"signatureCipher"`

	ContentLength int64  `json:"contentLength,string"`
	QualityLabel  string `json:"qualityLabel"`
}
