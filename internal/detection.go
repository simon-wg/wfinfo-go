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
	item, _, _ := seekInRow(lower, trie.Root)
	if item != nil {
		return item, nil
	}

	// Try upper row
	item, node, consumed := seekInRow(upper, trie.Root)
	if item != nil {
		return item, nil
	}

	if consumed {
		// Partial match in upper, try to continue in lower
		item, _, _ = seekInRow(lower, node)
		if item != nil {
			return item, nil
		}
	}

	return nil, fmt.Errorf("no item found")
}

func seekInRow(row *[]string, startNode *TrieNode) (*wfm.Item, *TrieNode, bool) {
	ptr := 0
	currNode := startNode
	wordsFound := 0

	for ptr < len(*row) {
		word := (*row)[ptr]
		if nextNode, ok := currNode.Children[word]; ok {
			currNode = nextNode
			wordsFound++
			ptr++
		} else {
			if wordsFound > 0 {
				// We were matching and it stopped.
				*row = append((*row)[:ptr-wordsFound], (*row)[ptr:]...)
				return currNode.Item, currNode, true
			}
			// If we are at the start (Root), we skip this word and keep looking for a match start
			if currNode == startNode {
				ptr++
			} else {
				// We were trying to continue a match from a previous row, but it failed immediately.
				return nil, startNode, false
			}
		}
	}

	if wordsFound > 0 {
		*row = append((*row)[:ptr-wordsFound], (*row)[ptr:]...)
		return currNode.Item, currNode, true
	}

	return nil, startNode, false
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

	// Split the text into words
	for word := range strings.FieldsSeq(text) {
		word = strings.TrimSpace(word)
		if correctedWord, ok := closestLegalWord(word, legalWords, 2); ok {
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
