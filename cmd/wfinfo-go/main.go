package main

import (
	"github.com/anthonynsimon/bild/imgio"
	"github.com/simon-wg/wfinfo-go/internal"
)

func main() {
	// Placeholder until screencaps work
	img, err := imgio.Open("internal/testdata/vitruvian-1.png")
	if err != nil {
		panic(err)
	}

	items := internal.DetectItems(img)
	for _, item := range items {
		println("Found item:", item.I18N["en"].Name)
	}
}
