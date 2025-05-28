package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/simon-wg/wf-ocr/internal/wfm"
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

// getItemsFromFile retrieves items from a local cache file if it exists and is not older than 24 hours.
func getItemsFromFile() []wfm.ItemJson {
	file, err := os.ReadFile("items.json")
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
	c := http.Client{}

	// Fetch items from the Warframe Market API
	resp, err := c.Get(itemsUrl)
	if err != nil {
		panic(fmt.Errorf("Failed to fetch items: %w", err))
	}
	defer resp.Body.Close()

	// Handle rate limiting by retrying after a delay
	// If it fails 5 times, or has a different error, panic
	failureCount := 0
	if resp.StatusCode == http.StatusTooManyRequests && failureCount <= 5 {
		failureCount++
		sleepDuration := math.Pow(2, float64(failureCount))
		time.Sleep(time.Duration(sleepDuration) * time.Second)
		resp, err = c.Get(itemsUrl)
	} else if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("Failed to fetch items, status code: %d", resp.StatusCode))
	}

	// Decode the JSON response from the API
	var itemResponse wfm.ItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&itemResponse); err != nil {
		panic(fmt.Errorf("Failed to decode API response: %w", err))
	}

	// Reencode the data to ensure it is in the correct format
	rawItems, err := json.Marshal(itemResponse.Data)
	if err != nil {
		panic(fmt.Errorf("Failed to marshal items: %w", err))
	}

	// Unmarshal the raw items into a slice of ItemJson
	var items []wfm.ItemJson
	if err := json.Unmarshal(rawItems, &items); err != nil {
		panic(fmt.Errorf("Failed to unmarshal items: %w", err))
	}

	// Add a timestamp to the structure and write it to a file
	var itemStore ItemStore
	itemStore.Timestamp = time.Now()
	itemStore.Items = items
	jsonData, err := json.MarshalIndent(itemStore, "", "  ")
	if err != nil {
		panic(fmt.Errorf("Failed to marshal items to JSON: %w", err))
	}
	if err := os.WriteFile("items.json", jsonData, 0644); err != nil {
		panic(fmt.Errorf("Failed to write items to file: %w", err))
	}

	return items
}

type ItemStore struct {
	Timestamp time.Time      `json:"timestamp"`
	Items     []wfm.ItemJson `json:"items"`
}
