package internal_test

import (
	"testing"

	"github.com/disintegration/imaging"
	"github.com/simon-wg/wfinfo-go/internal"
)

func TestOCR(t *testing.T) {
	img, err := imaging.Open("testdata/legacy-1.png")
	if err != nil {
		t.Fatalf("Error opening image: %v", err)
	}
	items, err := internal.GetItemsFromImage(img)
	if err != nil {
		t.Fatalf("Error getting items from image: %v", err)
	}
	expectedItems := []string{"Octavia Prime Blueprint", "Tenora Prime Blueprint", "Octavia Prime Systems Blueprint", "Harrow Prime Systems Blueprint"}
	actualItems := []string{}

	for _, item := range items {
		if item.I18N["en"] != nil {
			actualItems = append(actualItems, item.I18N["en"].Name)
		}
	}
	if len(actualItems) != len(expectedItems) {
		t.Fatalf("Expected %d items, got %d", len(expectedItems), len(actualItems))
	}
	for _, expected := range expectedItems {
		found := false
		for _, actual := range actualItems {
			if expected == actual {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected item '%s' not found in actual items", expected)
		}
	}
	for _, item := range actualItems {
		found := false
		for _, expected := range expectedItems {
			if expected == item {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected item '%s' found in actual items", item)
		}
	}
}
