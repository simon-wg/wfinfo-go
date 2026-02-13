package internal

import (
	"slices"
	"strings"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func GetPrimeItems() []wfm.Item {
	return filterPrimeItems(GetItems())
}

func filterPrimeItems(items []wfm.Item) []wfm.Item {
	forma := wfm.Item{
		Id:      "forma",
		Slug:    "forma",
		GameRef: "forma",
		Tags:    []string{"forma"},
		I18N: map[string]*wfm.ItemI18N{
			"en": {
				Name: "Forma Blueprint",
			},
		},
	}

	var primeItems []wfm.Item
	// Filter items that can be found in relics
	for _, item := range items {
		if item.Tags == nil {
			continue // Skip items without tags
		}
		if !slices.Contains(item.Tags, "prime") || !(slices.Contains(item.Tags, "blueprint") || slices.Contains(item.Tags, "component")) {
			continue
		}
		primeItems = append(primeItems, item)
	}

	primeItems = append(primeItems, forma) // Add Forma to the items fetched from the API

	return primeItems
}

// GetItems retrieves all items from the Warframe Market API or from a local cache file.
func GetItems() []wfm.Item {
	client := wfm.NewClient()
	items, err := client.FetchItems()
	if err != nil {
		panic(err)
	}

	return items
}

// Provides a list of all legal words that can be used in item names.
func AllLegalWords() []string {
	return extractLegalWords(GetPrimeItems())
}

func extractLegalWords(items []wfm.Item) []string {
	legalWords := map[string]struct{}{}

	for _, item := range items {
		for word := range strings.SplitSeq(item.I18N["en"].Name, " ") {
			legalWords[word] = struct{}{} // Use a map to ensure uniqueness
		}
	}
	words := []string{"Forma"}
	for word := range legalWords {
		if word == "Forma" {
			continue
		}
		words = append(words, word)
	}
	return words
}
