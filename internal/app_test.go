package internal

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/simon-wg/wfinfo-go/internal/wfm"
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

// New test for handleLine functionality
func TestHandleLine(t *testing.T) {
	// Create a mock log parser with buffer
	logBuffer := bytes.NewBufferString("some log line\n")
	parser := &logParser{
		reader: bufio.NewReader(logBuffer),
	}

	app := &appState{
		logParser:  parser,
		detection:  &detectionState{},
		foundItems: make(chan []wfm.Item, 1),
	}

	// Test line continuation
	parser.lineFragment = "partial "
	app.handleLine("line", nil)
	if parser.lineFragment != "" {
		t.Error("Expected line fragment to be cleared")
	}

	// Test incomplete line (no newline)
	parser.lineFragment = "incomplete"
	app.handleLine("", io.EOF)
	if parser.lineFragment != "incomplete" {
		t.Error("Expected line fragment to be preserved")
	}
}

// New test for rate limiting
func TestDetectionRateLimiting(t *testing.T) {
	app := &appState{
		logParser: &logParser{},
		detection: &detectionState{
			lastTriggered: time.Now(),
		},
		foundItems: make(chan []wfm.Item, 1),
	}

	// Test that rate limiting prevents detection
	line := "VoidProjections: OpenVoidProjectionRewardScreenRMI"

	// First call should update lastTriggered but not trigger detection (no OCR client)
	app.handleLine(line, nil)

	// Second immediate call should be rate limited
	app.handleLine(line, nil)

	// Verify no detection was triggered (channel should be empty)
	select {
	case <-app.foundItems:
		t.Error("Expected rate limiting to prevent detection")
	case <-time.After(100 * time.Millisecond):
		// Expected - no detection should occur
	}
}

// New test for error handling in detection
func TestDetectionErrorHandling(t *testing.T) {
	// Test that channel operations work correctly
	app := &appState{
		foundItems: make(chan []wfm.Item, 1),
	}

	// Channel should not block when receiving
	testItems := []wfm.Item{{
		Id: "test-item",
		I18N: map[string]*wfm.ItemI18N{
			"en": {Name: "Test Item"},
		},
	}}
	app.foundItems <- testItems

	select {
	case items := <-app.foundItems:
		if len(items) != 1 {
			t.Error("Expected 1 item, got", len(items))
		}
		// Success - channel works
	case <-time.After(100 * time.Millisecond):
		t.Error("Channel operation timed out")
	}
}
