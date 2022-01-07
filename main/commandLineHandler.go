package main

import (
	"context"

	// For command line operations
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "mp4",
	Short:   "Downloads a video from youtube",
	Example: `./main download Jl8fV1jUQPs -> for downloading https://www.youtube.com/watch?v=Jl8fV1jUQPs`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exitOnError(download(args[0]))
	},
}

// Initializes rootCommand and waits for
func init() {
	rootCmd.AddCommand(downloadCmd)
}

// download is the highest level functionality for downloading, currently only works for mp4
func download(id string) error {
	video, format, err := getVideoWithFormat(id)
	if err != nil {
		return err
	}

	return downloader.Download(context.Background(), video, format)
}
