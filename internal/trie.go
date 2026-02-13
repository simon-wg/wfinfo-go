package internal

import (
	"strings"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

type TrieNode struct {
	Children map[string]*TrieNode
	Item     *wfm.Item
}

type Trie struct {
	Root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		Root: &TrieNode{
			Children: make(map[string]*TrieNode),
		},
	}
}

func (t *Trie) Insert(item wfm.Item) {
	name := item.I18N["en"].Name
	words := strings.Split(name, " ")
	node := t.Root
	for _, word := range words {
		if _, ok := node.Children[word]; !ok {
			node.Children[word] = &TrieNode{
				Children: make(map[string]*TrieNode),
			}
		}
		node = node.Children[word]
	}
	node.Item = &item
}

func BuildTrie(items []wfm.Item) *Trie {
	trie := NewTrie()
	for _, item := range items {
		trie.Insert(item)
	}
	return trie
}
