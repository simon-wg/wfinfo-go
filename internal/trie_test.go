package internal

import (
	"testing"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func TestTrie(t *testing.T) {
	items := []wfm.Item{
		{
			I18N: map[string]*wfm.ItemI18N{
				"en": {Name: "Chassis Prime"},
			},
		},
		{
			I18N: map[string]*wfm.ItemI18N{
				"en": {Name: "Systems Prime"},
			},
		},
		{
			I18N: map[string]*wfm.ItemI18N{
				"en": {Name: "Chassis"},
			},
		},
	}

	trie := BuildTrie(items)

	// Test "Chassis Prime"
	node := trie.Root
	for _, word := range []string{"Chassis", "Prime"} {
		nextNode, ok := node.Children[word]
		if !ok {
			t.Fatalf("Expected word %q not found in trie", word)
		}
		node = nextNode
	}
	if node.Item == nil || node.Item.I18N["en"].Name != "Chassis Prime" {
		t.Errorf("Expected item %q at node, got %v", "Chassis Prime", node.Item)
	}

	// Test "Systems Prime"
	node = trie.Root
	for _, word := range []string{"Systems", "Prime"} {
		nextNode, ok := node.Children[word]
		if !ok {
			t.Fatalf("Expected word %q not found in trie", word)
		}
		node = nextNode
	}
	if node.Item == nil || node.Item.I18N["en"].Name != "Systems Prime" {
		t.Errorf("Expected item %q at node, got %v", "Systems Prime", node.Item)
	}

	// Test "Chassis" (prefix of "Chassis Prime")
	node = trie.Root.Children["Chassis"]
	if node.Item == nil || node.Item.I18N["en"].Name != "Chassis" {
		t.Errorf("Expected item %q at node, got %v", "Chassis", node.Item)
	}
}
