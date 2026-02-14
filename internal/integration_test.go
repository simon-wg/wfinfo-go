package internal_test

import (
	"testing"

	"github.com/disintegration/imaging"
	"github.com/simon-wg/wfinfo-go/internal"
)

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
			items, err := internal.GetItemsFromImage(img)
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
