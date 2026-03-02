package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/otiai10/gosseract/v2"
	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func Run(filePath string, steamLibrary string) error {
	fullPath, err := resolveEEPath(filePath, steamLibrary)
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Printf("Error closing watcher: %v", err)
		}
	}()

	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	if _, err := file.Seek(0, 2); err != nil {
		return err
	}
	if err := watcher.Add(fullPath); err != nil {
		return err
	}

	ocrClient := gosseract.NewClient()
	if err := ocrClient.SetWhitelist("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ& \n"); err != nil {
		return fmt.Errorf("failed to configure OCR: %w", err)
	}

	s := &appState{
		logParser: &logParser{
			reader: bufio.NewReader(file),
		},
		detection:  &detectionState{},
		foundItems: make(chan []wfm.Item),
		ocrClient:  ocrClient,
	}
	defer func() {
		if err := ocrClient.Close(); err != nil {
			log.Printf("Error closing OCR client: %v", err)
		}
	}()

	log.Printf("Watching %s for relic screen\n", fullPath)

	wfmClient := wfm.NewClient()

	for {
		select {
		case items := <-s.foundItems:
			for _, item := range items {
				detailedInfo, err := wfmClient.FetchItemTopOrders(item.Id, nil)
				if err != nil {
					log.Printf("Error: Unable to fetch price information for %v, %v\n", item.Id, err)
				}
				var sumPrice float32
				for _, order := range detailedInfo.Sell {
					sumPrice += float32(order.Platinum)
				}
				// Ex. Tekko Prime Gauntlets - 2.75p, 20 ducats
				fmt.Printf("%v - %.2fp, %v ducats\n", item.I18N["en"].Name, sumPrice/float32(len(detailedInfo.Sell)), item.Ducats)
			}
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Write) {
				s.handleWriteEvent()
			}
		case err := <-watcher.Errors:
			return err
		}
	}
}

type detectionState struct {
	mu            sync.Mutex
	lastTriggered time.Time
}

type logParser struct {
	reader       *bufio.Reader
	lineFragment string
	mu           sync.Mutex
}

type appState struct {
	logParser  *logParser
	detection  *detectionState
	foundItems chan []wfm.Item
	ocrClient  *gosseract.Client
}

func (s *appState) handleWriteEvent() {
	for {
		line, err := s.logParser.reader.ReadString('\n')
		if line != "" {
			s.handleLine(line, err)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
	}
}

func (s *appState) handleLine(line string, err error) {
	line = s.logParser.lineFragment + line
	if err != nil {
		s.logParser.mu.Lock()
		s.logParser.lineFragment = line
		s.logParser.mu.Unlock()
		return
	}

	s.logParser.mu.Lock()
	s.logParser.lineFragment = ""
	s.logParser.mu.Unlock()

	if !processLogLine(strings.TrimSpace(line)) {
		return
	}

	s.detection.mu.Lock()
	defer s.detection.mu.Unlock()

	if time.Since(s.detection.lastTriggered) < 1*time.Minute {
		return
	}

	s.detection.lastTriggered = time.Now()
	go s.triggerDetection()
}

func (s *appState) triggerDetection() {
	time.Sleep(500 * time.Millisecond)
	img := screenshot()

	// img, _ := imgio.Open("internal/testdata/conquera-1.png")
	log.Println("detecting items")
	s.foundItems <- DetectItems(img, s.ocrClient)
}

func processLogLine(line string) bool {
	return strings.Contains(line, "VoidProjections: OpenVoidProjectionRewardScreenRMI") || strings.Contains(line, "ProjectionRewardChoice.lua: Relic rewards initialized") || strings.Contains(line, "VoidProjections: GetVoidProjectionRewards")
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	}
	return filepath.Clean(path), nil
}

func resolveEEPath(filePath, steamLibrary string) (string, error) {
	if filePath != "" {
		return expandPath(filePath)
	}

	path, err := expandPath(steamLibrary)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "steamapps/compatdata/230410/pfx/drive_c/users/steamuser/AppData/Local/Warframe/EE.log"), nil
}
