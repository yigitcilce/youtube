package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Youtube video downloader",
	Long:  `Created for Yemeksepeti Golang bootcamp`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hi TechCareer")
	},
}

// cobra guidelines for this configuration
func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.main.yaml)")
}
