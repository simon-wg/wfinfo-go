package internal

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"github.com/otiai10/gosseract/v2"
	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func DetectItems(img image.Image) []wfm.Item {
	// This only works for 1080p, 4 items
	// Good enough for the simple case
	// (px, py, dx, dy)
	rects := []image.Rectangle{
		image.Rect(477, 412, 477+239, 412+50),
		image.Rect(719, 412, 719+239, 412+50),
		image.Rect(962, 412, 962+239, 412+50),
		image.Rect(1204, 412, 1204+239, 412+50),
	}
	client := gosseract.NewClient()
	defer func() { _ = client.Close() }()
	if err := client.SetWhitelist("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ& \n"); err != nil {
		return nil
	}
	textColor := detectTextColor(&img)

	relicItems := getRelicItems()
	relicItemNames := getItemNames(relicItems)

	items := make([]wfm.Item, 0, len(rects))
	for _, rect := range rects {
		itemName, _ := detectItemInBox(&img, rect, client, textColor)
		if itemName == nil {
			continue
		}
		item := findBestItem(*itemName, relicItems, relicItemNames)
		items = append(items, item)
	}

	return items
}

func detectItemInBox(img *image.Image, rect image.Rectangle, client *gosseract.Client, textColor color.RGBA) (*string, error) {
	cropped := transform.Crop(*img, rect)
	isolated := isolateTargetColor(cropped, textColor, 60)
	imgBuf := new(bytes.Buffer)
	if err := png.Encode(imgBuf, isolated); err != nil {
		return nil, err
	}
	if err := client.SetImageFromBytes(imgBuf.Bytes()); err != nil {
		return nil, err
	}
	text, err := client.Text()
	if err != nil {
		return nil, err
	}
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	return &text, nil
}

func detectTextColor(img *image.Image) color.RGBA {
	rect := image.Rect(320, 52, 324, 82)
	sample := transform.Crop(*img, rect)
	var red, green, blue, alpha uint64
	pixels := sample.Pix
	for i := 0; i < len(pixels); i += 4 {
		red += uint64(pixels[i])
		green += uint64(pixels[i+1])
		blue += uint64(pixels[i+2])
		alpha += uint64(pixels[i+3])
	}
	pixelCount := uint64(len(pixels) / 4)
	averageColor := color.RGBA{
		R: uint8(red / pixelCount),
		G: uint8(green / pixelCount),
		B: uint8(blue / pixelCount),
		A: uint8(alpha / pixelCount),
	}
	return averageColor
}

func isolateTargetColor(img *image.RGBA, target color.RGBA, threshold float64) *image.RGBA {
	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)
	r0, g0, b0 := float64(target.R), float64(target.G), float64(target.B)

	for i := 0; i < len(img.Pix); i += 4 {
		r1 := float64(img.Pix[i])
		g1 := float64(img.Pix[i+1])
		b1 := float64(img.Pix[i+2])

		dr := r0 - r1
		dg := g0 - g1
		db := b0 - b1
		dist := math.Sqrt(dr*dr + dg*dg + db*db)

		if dist < threshold {
			dest.Pix[i] = 0
			dest.Pix[i+1] = 0
			dest.Pix[i+2] = 0
			dest.Pix[i+3] = 255
		} else {
			dest.Pix[i] = 255
			dest.Pix[i+1] = 255
			dest.Pix[i+2] = 255
			dest.Pix[i+3] = 255
		}
	}
	return dest
}
