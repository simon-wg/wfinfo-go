package internal

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/disintegration/imaging"
)

func loadTestImage(t *testing.T, path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("could not open test image: %v", err)
	}
	defer func() { _ = file.Close() }()

	img, err := png.Decode(file)
	if err != nil {
		t.Fatalf("could not decode test image: %v", err)
	}
	return img
}

func TestDetectTextColor(t *testing.T) {
	testCases := []struct {
		name          string
		imagePath     string
		expectedColor color.RGBA
	}{
		{"conquera", "testdata/conquera-1.png", color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{"contrast", "testdata/contrast-1.png", color.RGBA{R: 102, G: 176, B: 255, A: 255}},
		{"equinox", "testdata/equinox-1.png", color.RGBA{R: 158, G: 159, B: 167, A: 255}},
		{"harrier", "testdata/harrier-1.png", color.RGBA{R: 253, G: 132, B: 2, A: 255}},
		{"legacy", "testdata/legacy-1.png", color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{"renewal", "testdata/renewal-1.png", color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{"vitruvian", "testdata/vitruvian-1.png", color.RGBA{R: 190, G: 169, B: 102, A: 255}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img := loadTestImage(t, tc.imagePath)
			actualColor := detectTextColor(&img)
			if actualColor != tc.expectedColor {
				t.Errorf("expected color %v, but got %v", tc.expectedColor, actualColor)
			}
		})
	}
}

func TestIsolateTargetColor(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.RGBA{R: 0, G: 0, B: 255, A: 255})

	targetColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	threshold := 10.0

	isolatedImg := isolateTargetColor(img, targetColor, threshold)

	// The pixel that is the same color as the target should be black
	if isolatedImg.At(0, 0) != (color.RGBA{R: 0, G: 0, B: 0, A: 255}) {
		t.Errorf("Expected pixel at (0,0) to be black, but got %v", isolatedImg.At(0, 0))
	}

	// The pixel that is a different color should be white
	if isolatedImg.At(1, 0) != (color.RGBA{R: 255, G: 255, B: 255, A: 255}) {
		t.Errorf("Expected pixel at (1,0) to be white, but got %v", isolatedImg.At(1, 0))
	}
}

func TestOCR(t *testing.T) {
	tests := []struct {
		name          string
		imagePath     string
		expectedItems []string
	}{
		{
			name:      "Conquera",
			imagePath: "testdata/conquera-1.png",
			expectedItems: []string{
				"Masseter Prime Handle",
				"Epitaph Prime Barrel",
				"Titania Prime Systems Blueprint",
				"Trumna Prime Blueprint",
			},
		},
		{
			name:      "Contrast",
			imagePath: "testdata/contrast-1.png",
			expectedItems: []string{
				"Burston Prime Receiver",
				"Orthos Prime Handle",
				"Ash Prime Neuroptics Blueprint",
				"Sevagoth Prime Systems Blueprint",
			},
		},
		{
			name:      "Equinox",
			imagePath: "testdata/equinox-1.png",
			expectedItems: []string{
				"Dual Zoren Prime Handle",
				"Bronco Prime Blueprint",
				"Alternox Prime Barrel",
				"Trumna Prime Blueprint",
			},
		},
		{
			name:      "Harrier",
			imagePath: "testdata/harrier-1.png",
			expectedItems: []string{
				"Grendel Prime Chassis Blueprint",
				"Cernos Prime Grip",
				"Bo Prime Blueprint",
				"Quassus Prime Blueprint",
			},
		},
		{
			name:      "Legacy",
			imagePath: "testdata/legacy-1.png",
			expectedItems: []string{
				"Hildryn Prime Systems Blueprint",
				"Mesa Prime Blueprint",
				"Caliban Prime Chassis Blueprint",
				"Bronco Prime Blueprint",
			},
		},
		{
			name:      "Renewal",
			imagePath: "testdata/renewal-1.png",
			expectedItems: []string{
				"Daikyu Prime Blueprint",
				"Acceltra Prime Receiver",
				"Caliban Prime Chassis Blueprint",
				"Lavos Prime Chassis Blueprint",
			},
		},
		{
			name:      "Vitruvian",
			imagePath: "testdata/vitruvian-1.png",
			expectedItems: []string{
				"Octavia Prime Blueprint",
				"Tenora Prime Blueprint",
				"Octavia Prime Systems Blueprint",
				"Harrow Prime Systems Blueprint",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := imaging.Open(tt.imagePath)
			if err != nil {
				t.Fatalf("Error opening image: %v", err)
			}
			items := DetectItems(img)
			if err != nil {
				t.Fatalf("Error getting items from image: %v", err)
			}

			actualItems := make([]string, 0, len(items))
			for _, item := range items {
				if en, ok := item.I18N["en"]; ok {
					actualItems = append(actualItems, en.Name)
				}
			}

			if len(tt.expectedItems) == 0 {
				t.Logf("Actual items for %s: %v", tt.name, actualItems)
				return
			}

			if len(actualItems) != len(tt.expectedItems) {
				t.Errorf("Expected %d items, got %d. Actual: %v", len(tt.expectedItems), len(actualItems), actualItems)
				return
			}

			expectedMap := make(map[string]bool)
			for _, item := range tt.expectedItems {
				expectedMap[item] = true
			}

			actualMap := make(map[string]bool)
			for _, item := range actualItems {
				actualMap[item] = true
			}

			for item := range expectedMap {
				if !actualMap[item] {
					t.Errorf("Expected item '%s' not found", item)
				}
			}

			for item := range actualMap {
				if !expectedMap[item] {
					t.Errorf("Unexpected item '%s' found", item)
				}
			}
		})
	}
}
