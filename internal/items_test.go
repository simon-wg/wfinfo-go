package internal

import (
	"slices"
	"testing"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func TestFilterPrimeItems(t *testing.T) {
	items := []wfm.Item{
		{
			Tags: []string{"prime", "warframe", "blueprint"},
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Excalibur Prime Blueprint"}},
		},
		{
			Tags: []string{"prime", "warframe"}, // Should be filtered (missing blueprint)
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Excalibur Prime"}},
		},
		{
			Tags: []string{"prime", "weapon", "blueprint"},
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Braton Prime Blueprint"}},
		},
		{
			Tags: []string{"prime", "weapon", "component"},
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Braton Prime Barrel"}},
		},
		{
			Tags: []string{"prime", "weapon"}, // Should be filtered (missing bp/component)
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Braton Prime"}},
		},
		{
			Tags: []string{"not-prime"}, // Should be filtered
			I18N: map[string]*wfm.ItemI18N{"en": {Name: "Braton"}},
		},
	}

	filtered := filterPrimeItems(items)

	expectedNames := []string{
		"Excalibur Prime Blueprint",
		"Braton Prime Blueprint",
		"Braton Prime Barrel",
	}

	if len(filtered) != len(expectedNames) {
		t.Fatalf("expected %d items, got %d", len(expectedNames), len(filtered))
	}

	for _, name := range expectedNames {
		found := false
		for _, item := range filtered {
			if item.I18N["en"].Name == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected item %s not found in filtered list", name)
		}
	}
}

func TestGetItemNames(t *testing.T) {
	items := []wfm.Item{
		{I18N: map[string]*wfm.ItemI18N{"en": {Name: "Item 1"}}},
		{I18N: map[string]*wfm.ItemI18N{"en": {Name: "Item 2"}}},
	}

	names := getItemNames(items)
	expected := []string{"Item 1", "Item 2"}

	if !slices.Equal(names, expected) {
		t.Errorf("expected %v, got %v", expected, names)
	}
}
