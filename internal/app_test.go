package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessLogLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected bool
	}{
		{"reward screen line", "some other stuff VoidProjections: OpenVoidProjectionRewardScreenRMI and more", true},
		{"relic rewards initialized", "some other stuff ProjectionRewardChoice.lua: Relic rewards initialized and more", true},
		{"get void projection rewards", "some other stuff VoidProjections: GetVoidProjectionRewards and more", true},
		{"unrelated line", "this is a random log line", false},
		{"empty line", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := processLogLine(tc.line)
			if actual != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, actual)
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get user home directory: %v", err)
	}

	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{"normal path", "/foo/bar", "/foo/bar"},
		{"home path", "~/.config", filepath.Join(home, ".config")},
		{"empty path", "", "."},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := expandPath(tc.path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual != tc.expected {
				t.Errorf("expected %q, but got %q", tc.expected, actual)
			}
		})
	}
}

func TestResolveEEPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get user home directory: %v", err)
	}

	testCases := []struct {
		name         string
		filePath     string
		steamLibrary string
		expected     string
	}{
		{"with file path", "/foo/bar", "", "/foo/bar"},
		{"with steam library", "", "/steam", "/steam/steamapps/compatdata/230410/pfx/drive_c/users/steamuser/AppData/Local/Warframe/EE.log"},
		{"with home in steam library", "", "~/.steam", filepath.Join(home, ".steam/steamapps/compatdata/230410/pfx/drive_c/users/steamuser/AppData/Local/Warframe/EE.log")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := resolveEEPath(tc.filePath, tc.steamLibrary)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual != tc.expected {
				t.Errorf("expected %q, but got %q", tc.expected, actual)
			}
		})
	}
}
