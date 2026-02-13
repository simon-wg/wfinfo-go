package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func GetPrimeItems() []wfm.ItemJson {
	items := GetItems()

	forma := wfm.ItemJson{
		Id:      "forma",
		Slug:    "forma",
		GameRef: "forma",
		Tags:    []string{"forma"},
		I18N: map[string]*wfm.ItemI18NJson{
			"en": {
				Name: "Forma Blueprint",
			},
		},
	}

	var primeItems []wfm.ItemJson
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
func GetItems() []wfm.ItemJson {
	var items []wfm.ItemJson

	items = getItemsFromFile()
	if items != nil {
		return items
	}

	items = getItemsFromApi()

	return items
}

func GetItemByNameSlice(nameSlice []string) *wfm.ItemJson {
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

func getItemsPath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "items.json"
	}
	appCacheDir := filepath.Join(cacheDir, "wfinfo-go")
	_ = os.MkdirAll(appCacheDir, 0755)
	return filepath.Join(appCacheDir, "items.json")
}

// getItemsFromFile retrieves items from a local cache file if it exists and is not older than 24 hours.
func getItemsFromFile() []wfm.ItemJson {
	file, err := os.ReadFile(getItemsPath())
	if err != nil {
		return nil // File does not exist or cannot be read
	}
	itemStore := ItemStore{}
	json.Unmarshal(file, &itemStore)
	if itemStore.Timestamp.Add(24 * time.Hour).After(time.Now()) {
		return itemStore.Items // Return cached items if they are not older than 24 hours
	}
	return nil // Cache is outdated or not available
}

// getItemsFromApi fetches items from the Warframe Market API and caches them to a local file.
func getItemsFromApi() []wfm.ItemJson {
	items, err := wfm.FetchItems()
	if err != nil {
		panic(err)
	}

	// Add a timestamp to the structure and write it to a file
	var itemStore ItemStore
	itemStore.Timestamp = time.Now()
	itemStore.Items = items
	jsonData, err := json.MarshalIndent(itemStore, "", "  ")
	if err != nil {
		panic(fmt.Errorf("Failed to marshal items to JSON: %w", err))
	}
	if err := os.WriteFile(getItemsPath(), jsonData, 0644); err != nil {
		panic(fmt.Errorf("Failed to write items to file: %w", err))
	}

	return items
}

type ItemStore struct {
	Timestamp time.Time      `json:"timestamp"`
	Items     []wfm.ItemJson `json:"items"`
}
