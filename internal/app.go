package internal

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
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
	defer func() { _ = watcher.Close() }()

	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	if _, err := file.Seek(0, 2); err != nil {
		return err
	}
	if err := watcher.Add(fullPath); err != nil {
		return err
	}

	s := &state{
		reader:     bufio.NewReader(file),
		foundItems: make(chan []wfm.Item),
	}

	log.Printf("Watching %s for relic screen\n", fullPath)

	for {
		select {
		case items := <-s.foundItems:
			for _, item := range items {
				log.Println(item.I18N["en"].Name)
			}
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Write) {
				s.handleWriteEvent()
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

type state struct {
	reader        *bufio.Reader
	lineFragment  string
	lastTriggered time.Time
	foundItems    chan []wfm.Item
}

func (s *state) handleWriteEvent() {
	for {
		line, err := s.reader.ReadString('\n')
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

func (s *state) handleLine(line string, err error) {
	line = s.lineFragment + line
	if err != nil {
		s.lineFragment = line
		return
	}

	s.lineFragment = ""
	if !processLogLine(strings.TrimSpace(line)) {
		return
	}

	if time.Since(s.lastTriggered) < 1*time.Minute {
		return
	}

	s.lastTriggered = time.Now()
	go s.triggerDetection()
}

func (s *state) triggerDetection() {
	time.Sleep(500 * time.Millisecond)
	img := screenshot()

	// img, _ := imgio.Open("internal/testdata/conquera-1.png")
	log.Println("detecting items")
	s.foundItems <- DetectItems(img)
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
