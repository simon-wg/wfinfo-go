package internal

import (
	"bytes"
	"fmt"
	"image"
	"slices"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"github.com/simon-wg/wf-ocr/internal/wfm"
)

func GetItemsFromImage(img image.Image) ([]wfm.ItemJson, error) {
	client := gosseract.NewClient()
	client.SetWhitelist("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ")
	defer client.Close()

	img = preprocessImage(img)
	items := getItemsFromImage(client, img)

	return items, nil
}

func getItemsFromImage(c *gosseract.Client, img image.Image) []wfm.ItemJson {
	upperImg := imaging.Crop(img, image.Rect(0, 0, 965, 25))
	lowerImg := imaging.Crop(img, image.Rect(0, 25, 965, 50))

	upperWords, err := getWordsFromImage(c, upperImg)
	if err != nil {
		panic(err)
	}
	lowerWords, err := getWordsFromImage(c, lowerImg)
	if err != nil {
		panic(err)
	}

	items := getItemsFromWords(&upperWords, &lowerWords)
	return items
}

func getItemsFromWords(upper *[]string, lower *[]string) []wfm.ItemJson {
	var foundItems []wfm.ItemJson
	allItems := GetPrimeItems()

	// What is done here is to first seek through the bottom row.
	// If a word which can be an item is found we check as far as possible on the
	// same row for the rest of the item name.
	// Once the full name is found we remove all words that were used from the rows.
	// If we can't find a full item name on the bottom row we check the top row for the first word/words
	// and then check the bottom row for the rest of the item name.
	for len(*upper) > 0 && len(*lower) > 0 {
		item, err := findItemInWords(upper, lower, allItems)
		if err == nil {
			foundItems = append(foundItems, *item)
		}
	}

	return foundItems
}

func findItemInWords(upper *[]string, lower *[]string, allItems []wfm.ItemJson) (*wfm.ItemJson, error) {
	candidate := []string{}
	// Seek through the lower row first
	candidate, wordIdx := seekItemName(candidate, lower, 0, allItems)
	if item := GetItemByNameSlice(candidate); item != nil {
		return item, nil
	}

	// If we reach here, we didn't find a valid item in the lower row
	// Now we check the upper row for the first word/words
	candidate = []string{}
	candidate, wordIdx = seekItemName(candidate, upper, 0, allItems)
	if item := GetItemByNameSlice(candidate); item != nil {
		return item, nil
	}

	// We found an incomplete item name in the upper row
	// Now we check the lower row for the rest of the item name
	candidate, wordIdx = seekItemName(candidate, lower, wordIdx, allItems)
	if item := GetItemByNameSlice(candidate); item != nil {
		return item, nil
	}

	return nil, fmt.Errorf("No item found") // Return an empty item if no match is found
}

// This function seeks a certain row for an item name.
// It returns the candidate name, the index of the word and the ptr of the row.
// A ptrIdx != -1 indicates that the item was found and the candidate is complete.
// A ptrIdx >= 0 indicates that the item was not found and the candidate is incomplete.
func seekItemName(candidate []string, row *[]string, wordIdx int, allItems []wfm.ItemJson) ([]string, int) {
	ptr := 0
	wordsFound := 0

	for ptr < len(*row) {
		word := (*row)[ptr]
		if wordIndexCorrect(candidate, word, wordIdx, allItems) {
			wordsFound++
			candidate = append(candidate, word)
			wordIdx++
			ptr++
		} else {
			if wordsFound > 0 {
				*row = append((*row)[:ptr-wordsFound], (*row)[ptr:]...) // Remove the words that were used from the row
				return candidate, wordIdx                               // Return the candidate if we have found some words
			} else {
				// If we haven't found any words, we need to skip this word
				ptr++
			}
		}
	}

	*row = append((*row)[:ptr-wordsFound], (*row)[ptr:]...) // Remove the words that were used from the row
	return candidate, wordIdx
}

// This checks if the word at the given index in the item name matches the word in the item list.
func wordIndexCorrect(candidate []string, word string, index int, allItems []wfm.ItemJson) bool {
	for _, item := range allItems {
		if item.I18N["en"] == nil {
			continue // Skip items without English localization
		}
		itemNameSplit := strings.Split(item.I18N["en"].Name, " ")
		if index >= len(itemNameSplit) {
			continue // Skip if the index is out of bounds for the item name
		}
		// Check if the candidate matches the item name up to the current index
		if slices.Equal(candidate, itemNameSplit[:index]) && itemNameSplit[index] == word {
			return true
		}
	}
	return false
}

func preprocessImage(img image.Image) image.Image {
	// Adjust contrast, invert colors, and convert to grayscale
	img = imaging.AdjustContrast(img, 80)
	img = imaging.Invert(img)
	img = imaging.Grayscale(img)

	// Crop the image to the relevant area
	cropRegion := image.Rect(480, 410, 480+965, 410+50)
	return imaging.Crop(img, cropRegion)
}

func getWordsFromImage(c *gosseract.Client, img image.Image) ([]string, error) {
	buffer := bytes.NewBuffer(nil)
	if err := imaging.Encode(buffer, img, imaging.PNG); err != nil {
		return nil, err
	}

	c.SetImageFromBytes(buffer.Bytes())

	// Perform OCR on the image
	text, err := c.Text()
	if err != nil {
		return nil, err
	}

	legalWords := AllLegalWords()
	words := []string{}

	// Split the text into words
	for _, word := range strings.Fields(text) {
		word = strings.TrimSpace(word)
		if legalWord(word, legalWords) {
			words = append(words, word)
		}
	}

	return words, nil
}

func legalWord(word string, legalWords []string) bool {
	if word == "forma" || word == "Forma" {
		return true // Forma is a special case that should always be considered legal
	}
	for _, legalWord := range legalWords {
		if strings.EqualFold(legalWord, word) {
			return true
		}
	}
	return false
}
