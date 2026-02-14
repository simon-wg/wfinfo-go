package internal

import (
	"bytes"
	"fmt"
	"image"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func GetItemsFromImage(img image.Image) ([]wfm.Item, error) {
	client := gosseract.NewClient()
	//nolint:errcheck
	defer client.Close()
	err := client.SetWhitelist("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ")
	if err != nil {
		return nil, err
	}
	img = preprocessImage(img)
	items := getItemsFromImage(client, img)
	return items, nil
}

func getItemsFromImage(c *gosseract.Client, img image.Image) []wfm.Item {
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

	allItems := GetPrimeItems()
	trie := BuildTrie(allItems)
	items := getItemsFromWords(&upperWords, &lowerWords, trie)
	return items
}

func getItemsFromWords(upper *[]string, lower *[]string, trie *Trie) []wfm.Item {
	var foundItems []wfm.Item

	for len(*upper) > 0 || len(*lower) > 0 {
		item, err := findItemInWords(upper, lower, trie)
		if err == nil {
			foundItems = append(foundItems, *item)
		} else {
			// If no item is found, consume one word to avoid infinite loop
			if len(*upper) > 0 {
				*upper = (*upper)[1:]
			} else if len(*lower) > 0 {
				*lower = (*lower)[1:]
			} else {
				break
			}
		}
	}

	return foundItems
}

func findItemInWords(upper *[]string, lower *[]string, trie *Trie) (*wfm.Item, error) {
	// Try lower row first
	item, matchCount, skipCount := findLongestMatchInRow(*lower, trie.Root)
	if item != nil {
		*lower = append((*lower)[:skipCount], (*lower)[skipCount+matchCount:]...)
		return item, nil
	}

	// Try upper row
	item, matchCount, skipCount = findLongestMatchInRow(*upper, trie.Root)
	if item != nil {
		*upper = append((*upper)[:skipCount], (*upper)[skipCount+matchCount:]...)
		return item, nil
	}

	// Try partial match in upper, continuing in lower
	// We only do this if we find a partial match at the start of the upper row (after skipping)
	// that doesn't form a full item on its own in the upper row.
	bestUpperItem, upperMatchCount, upperSkipCount, partialNode := seekInRow(*upper, trie.Root)
	if partialNode != nil && bestUpperItem == nil {
		item2, lowerMatchCount, lowerSkipCount, _ := seekInRow(*lower, partialNode)
		if item2 != nil {
			*upper = append((*upper)[:upperSkipCount], (*upper)[upperSkipCount+upperMatchCount:]...)
			*lower = append((*lower)[:lowerSkipCount], (*lower)[lowerSkipCount+lowerMatchCount:]...)
			return item2, nil
		}
	}

	return nil, fmt.Errorf("no item found")
}

func findLongestMatchInRow(row []string, startNode *TrieNode) (*wfm.Item, int, int) {
	item, matchCount, skipCount, _ := seekInRow(row, startNode)
	return item, matchCount, skipCount
}

func seekInRow(row []string, startNode *TrieNode) (*wfm.Item, int, int, *TrieNode) {
	bestItem := (*wfm.Item)(nil)
	bestMatchCount := 0
	bestSkipCount := 0
	var bestPartialNode *TrieNode

	for skip := 0; skip < len(row); skip++ {
		currNode := startNode
		matchCount := 0
		for i := skip; i < len(row); i++ {
			if nextNode, ok := currNode.Children[row[i]]; ok {
				currNode = nextNode
				matchCount++
				if currNode.Item != nil {
					bestItem = currNode.Item
					bestMatchCount = matchCount
					bestSkipCount = skip
					bestPartialNode = nil // Full match found, reset partial node
				} else if bestItem == nil {
					// Only track partial match if we haven't found a full one yet
					bestMatchCount = matchCount
					bestSkipCount = skip
					bestPartialNode = currNode
				}
			} else {
				break
			}
		}
		if bestItem != nil {
			return bestItem, bestMatchCount, bestSkipCount, nil
		}
		if bestPartialNode != nil {
			return nil, bestMatchCount, bestSkipCount, bestPartialNode
		}
	}

	return nil, 0, 0, nil
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

	err := c.SetImageFromBytes(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	// Perform OCR on the image
	text, err := c.Text()
	if err != nil {
		return nil, err
	}

	legalWords := AllLegalWords()
	words := []string{}

	// Split the text into words by whitespace, non-alphabetic characters, and case changes
	var currentWord strings.Builder
	var wordsToProcess []string
	for i, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			if i > 0 {
				prev := rune(text[i-1])
				// Split on lower -> UPPER
				if (prev >= 'a' && prev <= 'z') && (r >= 'A' && r <= 'Z') {
					if currentWord.Len() > 0 {
						wordsToProcess = append(wordsToProcess, currentWord.String())
						currentWord.Reset()
					}
				}
			}
			currentWord.WriteRune(r)
		} else {
			if currentWord.Len() > 0 {
				wordsToProcess = append(wordsToProcess, currentWord.String())
				currentWord.Reset()
			}
		}
	}
	if currentWord.Len() > 0 {
		wordsToProcess = append(wordsToProcess, currentWord.String())
	}

	for _, word := range wordsToProcess {
		word = strings.TrimSpace(word)
		if correctedWord, ok := closestLegalWord(word, legalWords, 3); ok {
			words = append(words, correctedWord)
		}
	}

	return words, nil
}

func closestLegalWord(word string, legalWords []string, maxDistance int) (string, bool) {
	bestWord := ""
	minDist := maxDistance + 1

	for _, legal := range legalWords {
		dist := levenshtein(strings.ToLower(word), strings.ToLower(legal))
		if dist < minDist {
			minDist = dist
			bestWord = legal
		}
	}

	if minDist <= maxDistance {
		return bestWord, true
	}

	return "", false
}
