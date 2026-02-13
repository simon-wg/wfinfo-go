package internal

import (
	"testing"
)

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		s, t string
		want int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"abc", "ab", 1},
		{"abc", "bc", 1},
		{"abc", "ac", 1},
		{"abc", "abcd", 1},
		{"kitten", "sitting", 3},
		{"forma", "format", 1},
		{"forma", "form", 1},
		{"forma", "f0rmal", 2},
	}

	for _, tt := range tests {
		got := levenshtein(tt.s, tt.t)
		if got != tt.want {
			t.Errorf("levenshtein(%q, %q) = %d; want %d", tt.s, tt.t, got, tt.want)
		}
	}
}
