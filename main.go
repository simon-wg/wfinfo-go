package main

import (
	"github.com/disintegration/imaging"
	"github.com/simon-wg/wfinfo-go/internal"
)

func main() {
	image, err := imaging.Open("test/test-images/1.png")
	if err != nil {
		panic(err)
	}

	items, err := internal.GetItemsFromImage(image)
	if err != nil {
		println("Error:", err)
	}
	for _, item := range items {
		println("Found item:", item.I18N["en"].Name)
	}
}
