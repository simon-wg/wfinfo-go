package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/fsnotify/fsnotify"
	"github.com/simon-wg/wfinfo-go/internal"
	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func main() {
	// filePath := flag.String("f", "~/.local/share/Steam/steamapps/compatdata/230410/pfx/drive_c/users/steamusers/AppData/Local/Warframe/EE.log", "Path to EE.log")
	filePath := flag.String("f", "~/Programming/Hobby/wfinfo-go/EE.log", "Path to EE.log")
	flag.Parse()
	fullPath, err := expandPath(*filePath)
	if err != nil {
		panic(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer func() { _ = watcher.Close() }()

	file, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()
	if _, err := file.Seek(0, 2); err != nil {
		panic(err)
	}
	if err := watcher.Add(fullPath); err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	log.Printf("Watching %s for relic screen\n", *filePath)
	foundItems := make(chan []wfm.Item)
	var lineFragment string
	var lastTriggered time.Time

	for {
		select {
		case items := <-foundItems:
			for _, item := range items {
				println(item.I18N["en"].Name)
			}
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				for {
					line, err := reader.ReadString('\n')
					if line != "" {
						line = lineFragment + line
						if err == nil {
							lineFragment = ""
							if !processLogLine(strings.TrimSpace(line)) {
								continue
							}
							if time.Since(lastTriggered) < 1*time.Minute {
								continue
							}
							lastTriggered = time.Now()
							go func() {
								img, _ := imgio.Open("internal/testdata/conquera-1.png")
								println("detecting items")
								foundItems <- internal.DetectItems(img)
							}()
						} else {
							lineFragment = line
						}
					}
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Printf("read error: %v", err)
					}
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
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

func processLogLine(line string) bool {
	return line == "test123"
}
