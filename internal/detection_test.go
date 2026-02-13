package internal

import (
	"testing"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func TestClosestLegalWord(t *testing.T) {
	legalWords := []string{"Prime", "Blueprint", "Systems", "Forma", "Harrow"}

	tests := []struct {
		word string
		want string
		ok   bool
	}{
		{"Prime", "Prime", true},
		{"Pirme", "Prime", true},
		{"Primey", "Prime", true},
		{"Pr", "", false}, // Distance 3
		{"forma", "Forma", true},
		{"f0rmal", "Forma", true},
		{"format", "Forma", true},
		{"Harrow", "Harrow", true},
		{"Harrowing", "", false}, // Distance 3
	}

	for _, tt := range tests {
		got, ok := closestLegalWord(tt.word, legalWords, 2)
		if ok != tt.ok {
			t.Errorf("closestLegalWord(%q) ok = %v; want %v", tt.word, ok, tt.ok)
		}
		if ok && got != tt.want {
			t.Errorf("closestLegalWord(%q) = %q; want %q", tt.word, got, tt.want)
		}
	}
}

func TestSeekInRow(t *testing.T) {
	item1 := wfm.Item{I18N: map[string]*wfm.ItemI18N{"en": {Name: "Ash Prime Chassi"}}}
	item2 := wfm.Item{I18N: map[string]*wfm.ItemI18N{"en": {Name: "Ash Prime Blueprint"}}}

	trie := NewTrie()
	trie.Insert(item1)
	trie.Insert(item2)

	tests := []struct {
		name          string
		row           []string
		expectedItem  string
		expectedWords int
		remainingRow  []string
	}{
		{
			name:          "Match full item",
			row:           []string{"Ash", "Prime", "Blueprint", "Other"},
			expectedItem:  "Ash Prime Blueprint",
			expectedWords: 3,
			remainingRow:  []string{"Other"},
		},
		{
			name:          "Match partial item (prefix)",
			row:           []string{"Ash", "Prime", "Chassi", "Other"},
			expectedItem:  "Ash Prime Chassi",
			expectedWords: 3,
			remainingRow:  []string{"Other"},
		},
		{
			name:          "No match at start",
			row:           []string{"Other", "Ash", "Prime", "Blueprint"},
			expectedItem:  "Ash Prime Blueprint",
			expectedWords: 3,
			remainingRow:  []string{"Other"}, // "Other" is skipped because it's at start and doesn't match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := tt.row
			item, _, _ := seekInRow(&row, trie.Root)
			if tt.expectedItem == "" {
				if item != nil {
					t.Errorf("Expected no item, got %v", item.I18N["en"].Name)
				}
			} else {
				if item == nil {
					t.Fatalf("Expected item %q, got nil", tt.expectedItem)
				}
				if item.I18N["en"].Name != tt.expectedItem {
					t.Errorf("Expected item %q, got %q", tt.expectedItem, item.I18N["en"].Name)
				}
			}
		})
	}
}

func TestGetItemsFromWords(t *testing.T) {
	item1 := wfm.Item{I18N: map[string]*wfm.ItemI18N{"en": {Name: "Ash Prime Blueprint"}}}
	item2 := wfm.Item{I18N: map[string]*wfm.ItemI18N{"en": {Name: "Tenora Prime Stock"}}}

	trie := NewTrie()
	trie.Insert(item1)
	trie.Insert(item2)

	tests := []struct {
		name          string
		upper         []string
		lower         []string
		expectedItems []string
	}{
		{
			name:          "Split item across rows",
			upper:         []string{"Ash", "Prime"},
			lower:         []string{"Blueprint"},
			expectedItems: []string{"Ash Prime Blueprint"},
		},
		{
			name:          "Multiple items (lower first)",
			upper:         []string{"Ash", "Prime", "Blueprint"},
			lower:         []string{"Tenora", "Prime", "Stock"},
			expectedItems: []string{"Tenora Prime Stock", "Ash Prime Blueprint"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upper := tt.upper
			lower := tt.lower
			items := getItemsFromWords(&upper, &lower, trie)
			if len(items) != len(tt.expectedItems) {
				t.Fatalf("Expected %d items, got %d", len(tt.expectedItems), len(items))
			}
			for i, item := range items {
				if item.I18N["en"].Name != tt.expectedItems[i] {
					t.Errorf("Expected item %d to be %q, got %q", i, tt.expectedItems[i], item.I18N["en"].Name)
				}
			}
		})
	}
}
