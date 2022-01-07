package main

import (
	"context"
	"io"
	"log"
	"os"
	"regexp"

	// For visualizing download progress
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"

	"github.com/yigitcilce/youtube"
)

// Downloader offers high level functions to download videos into files
type Downloader struct {
	youtube.Client
}

// DLProgress keeps track of downloaded content
type DLprogress struct {
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
}

// Write is an io.Wroter for download process
func (prog *DLprogress) Write(p []byte) (n int, err error) {
	n = len(p)
	prog.totalWrittenBytes = prog.totalWrittenBytes + float64(n)
	currentPercent := (prog.totalWrittenBytes / prog.contentLength) * 100
	if (prog.downloadLevel <= currentPercent) && (prog.downloadLevel < 100) {
		prog.downloadLevel++
	}
	return
}

// Download : Starting download video by arguments.
func (yt *Downloader) Download(ctx context.Context, v *youtube.Video, format *youtube.Format) error {

	// Tell me, what are you downloading
	yt.logf("Video '%s' - Quality '%s' - Codec '%s'", v.Title, format.QualityLabel, format.MimeType)

	// Create the file with video name and extension. only mp4 for now
	destFile := ConvertVideoTitletoFileName(v.Title) + ".mp4"

	// Create output file
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Go to real-deal, downloading process
	return yt.videoDLWorker(ctx, out, v, format)
}

// videoDLWorker starts the downloading process and visualize it to user in command line
func (yt *Downloader) videoDLWorker(ctx context.Context, out *os.File, video *youtube.Video, format *youtube.Format) error {
	stream, size, err := yt.GetStreamContext(ctx, video, format)
	if err != nil {
		return err
	}

	prog := &DLprogress{
		contentLength: float64(size),
	}

	// Configuration of the progress bar
	progress := mpb.New(mpb.WithWidth(100))
	bar := progress.AddBar(
		int64(prog.contentLength),

		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)

	reader := bar.ProxyReader(stream)
	mw := io.MultiWriter(out, prog)
	_, err = io.Copy(mw, reader)
	if err != nil {
		return err
	}

	progress.Wait()
	return nil
}

func (yt *Downloader) logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// ConvertVideoTitletoFileName gets rid of characters that are not allowed in operating systems
func ConvertVideoTitletoFileName(fileName string) string {
	// Not allowed on windows = <>:"/\|?* Mac = :/
	fileName = regexp.MustCompile(`[:/<>\:"\\|?*]`).ReplaceAllString(fileName, "")
	fileName = regexp.MustCompile(`\s+`).ReplaceAllString(fileName, " ")

	return fileName
}
