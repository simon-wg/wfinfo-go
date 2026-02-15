package main

import (
	"flag"

	"github.com/simon-wg/wfinfo-go/internal"
)

func main() {
	filePath := flag.String("f", "~/Programming/Hobby/wfinfo-go/EE.log", "Path to EE.log")
	flag.Parse()

	if err := internal.Run(*filePath); err != nil {
		panic(err)
	}
}
