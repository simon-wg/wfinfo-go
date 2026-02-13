package internal

import (
	"slices"
	"strings"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func GetPrimeItems() []wfm.Item {
	items := GetItems()

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
		if !slices.Contains(item.Tags, "prime") || slices.Contains(item.Tags, "set") {
			continue // Skip items that are not prime
		}
		if slices.Contains(item.Tags, "warframe") || slices.Contains(item.Tags, "weapon") || slices.Contains(item.Tags, "archwing") || slices.Contains(item.Tags, "sentinel") {
			primeItems = append(primeItems, item)
		}
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

func GetItemByNameSlice(nameSlice []string) *wfm.Item {
	name := strings.Join(nameSlice, " ")
	for _, item := range GetItems() {
		if item.I18N["en"].Name == name {
			return &item
		}
	}
	return nil
}

// Provides a list of all legal words that can be used in item names.
func AllLegalWords() []string {
	legalWords := map[string]struct{}{}

	for _, item := range GetPrimeItems() {
		for word := range strings.SplitSeq(item.I18N["en"].Name, " ") {
			legalWords[word] = struct{}{} // Use a map to ensure uniqueness
		}
	}
	words := []string{}
	for word := range legalWords {
		words = append(words, word)
	}
	return words
}
