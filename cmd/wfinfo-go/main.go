package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/simon-wg/wfinfo-go/internal"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	filePath := flag.String("f", "", "Path to EE.log (overrides -d)")
	steamLibrary := flag.String("d", "~/.local/share/Steam", "Path to Steam library folder")
	flag.Parse()

	if err := internal.Run(*filePath, *steamLibrary); err != nil {
		handleError(err, *filePath, *steamLibrary)
	}
}

func handleError(err error, filePath, steamLibrary string) {
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Invalid path to EE.log or Steam library.\n")
		if filePath != "" {
			fmt.Fprintf(os.Stderr, "Path: %s\n", filePath)
		} else {
			fmt.Fprintf(os.Stderr, "Steam library: %s\n", steamLibrary)
		}
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
	os.Exit(1)
}
