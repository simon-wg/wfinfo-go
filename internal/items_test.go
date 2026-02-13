package internal

import (
	"slices"
	"testing"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func TestFilterPrimeItems(t *testing.T) {
	items := []wfm.Item{
		{
			Tags: []string{"prime", "warframe"},
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Ash Prime"}},
		},
		{
			Tags: []string{"prime", "set", "warframe"},
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Ash Prime Set"}},
		},
		{
			Tags: []string{"weapon"},
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Braton"}},
		},
		{
			Tags: []string{"prime", "weapon"},
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Braton Prime"}},
		},
	}

	primeItems := filterPrimeItems(items)

	// Should contain Ash Prime, Braton Prime, and Forma Blueprint
	if len(primeItems) != 3 {
		t.Errorf("Expected 3 prime items, got %d", len(primeItems))
	}

	names := []string{}
	for _, item := range primeItems {
		names = append(names, item.I18N["en"].Name)
	}

	if !slices.Contains(names, "Ash Prime") {
		t.Errorf("Expected Ash Prime in prime items")
	}
	if !slices.Contains(names, "Braton Prime") {
		t.Errorf("Expected Braton Prime in prime items")
	}
	if !slices.Contains(names, "Forma Blueprint") {
		t.Errorf("Expected Forma Blueprint in prime items")
	}
	if slices.Contains(names, "Ash Prime Set") {
		t.Errorf("Did not expect Ash Prime Set in prime items")
	}
	if slices.Contains(names, "Braton") {
		t.Errorf("Did not expect Braton in prime items")
	}
}

func TestExtractLegalWords(t *testing.T) {
	items := []wfm.Item{
		{
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Ash Prime"}},
		},
		{
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Braton Prime"}},
		},
	}

	words := extractLegalWords(items)

	expectedWords := []string{"Forma", "Ash", "Prime", "Braton"}
	if len(words) != len(expectedWords) {
		t.Errorf("Expected %d words, got %d: %v", len(expectedWords), len(words), words)
	}

	for _, w := range expectedWords {
		if !slices.Contains(words, w) {
			t.Errorf("Expected word %q not found in %v", w, words)
		}
	}
}
