package main

import (
	"fmt"
	"os"
)

// If cmd returns an error, terminate the process
func main() {
	exitOnError(rootCmd.Execute())
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
