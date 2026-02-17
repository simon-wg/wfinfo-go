package main

import (
	"flag"

	"github.com/simon-wg/wfinfo-go/internal"
)

func main() {
	filePath := flag.String("f", "", "Path to EE.log")
	steamLibrary := flag.String("d", "~/.local/share/Steam", "Path to Steam library folder (should contain steamapps)")
	flag.Parse()

	if err := internal.Run(*filePath, *steamLibrary); err != nil {
		panic(err)
	}
}
