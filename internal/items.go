package internal

import (
	"slices"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func getRelicItems() []wfm.Item {
	client := wfm.NewClient()
	items, err := client.FetchItems()
	if err != nil {
		return nil
	}
	primes := filterPrimeItems(items)
	relicItems := append(primes, wfm.Item{
		Id:      "forma",
		Slug:    "forma",
		GameRef: "forma",
		Tags:    []string{"forma"},
		I18N: map[string]*wfm.ItemI18N{
			"en": {
				Name: "Forma Blueprint",
			},
		},
	})
	return relicItems
}

func filterPrimeItems(items []wfm.Item) []wfm.Item {
	primeItems := []wfm.Item{}
	for _, item := range items {
		// Make sure it's prime
		if !slices.Contains(item.Tags, "prime") {
			continue
		}
		// Make sure it's a weapon component/bp
		if slices.Contains(item.Tags, "weapon") && !(slices.Contains(item.Tags, "blueprint") || slices.Contains(item.Tags, "component")) {
			continue
		}
		// Make sure it's a valid warframe bp
		if slices.Contains(item.Tags, "warframe") && !(slices.Contains(item.Tags, "blueprint")) {
			continue
		}
		primeItems = append(primeItems, item)
	}
	return primeItems
}

func getItemNames(items []wfm.Item) []string {
	names := []string{}
	for _, item := range items {
		names = append(names, item.I18N["en"].Name)
	}
	return names
}

func findBestItem(itemName string, relicItems []wfm.Item, relicItemNames []string) wfm.Item {
	bestName := smithWaterman(itemName, relicItemNames)
	item := getItemFromName(bestName, relicItems)
	return item
}

func getItemFromName(name string, items []wfm.Item) wfm.Item {
	for _, item := range items {
		if item.I18N["en"].Name == name {
			return item
		}
	}
	panic("no item found")
}
