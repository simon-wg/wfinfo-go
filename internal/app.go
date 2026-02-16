package internal

import (
	"bufio"
	"image"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/fsnotify/fsnotify"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/simon-wg/wfinfo-go/internal/wfm"
)

func Run(filePath string) error {
	fullPath, err := expandPath(filePath)
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

	log.Printf("Watching %s for relic screen\n", filePath)

	for {
		select {
		case items := <-s.foundItems:
			for _, item := range items {
				println(item.I18N["en"].Name)
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
	// Screenshot instead
	img := s.screenshot()

	// img, _ := imgio.Open("internal/testdata/conquera-1.png")
	println("detecting items")
	s.foundItems <- DetectItems(img)
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

func (s *state) screenshot() image.Image {
	X, err := xgb.NewConn()
	if err != nil {
		log.Fatal("unable to find xserver")
	}
	defer X.Close()
	setup := xproto.Setup(X)
	root := setup.DefaultScreen(X).Root
	targetWin := findWindow(X, root, "steam_app_230410")
	if targetWin == 0 {
		log.Fatalf("unable to find Warframe window")
	}
	log.Printf("found warframe window id %d\n", targetWin)
	geom, err := xproto.GetGeometry(X, xproto.Drawable(targetWin)).Reply()
	if err != nil {
		log.Fatal(err)
	}
	reply, err := xproto.GetImage(X, xproto.ImageFormatZPixmap, xproto.Drawable(targetWin), 0, 0, geom.Width, geom.Height, 0xffffffff).Reply()
	if err != nil {
		log.Fatal(err)
	}
	img := x11ToImage(reply.Data, int(geom.Width), int(geom.Height))
	// INFO: This is only for debugging. Remove once tested on real window.
	_ = imgio.Save("screencap.png", img, imgio.PNGEncoder())
	return img
}

func findWindow(X *xgb.Conn, parent xproto.Window, targetClass string) xproto.Window {
	tree, err := xproto.QueryTree(X, parent).Reply()
	if err != nil {
		return 0
	}

	for _, child := range tree.Children {
		match := false
		prop, err := xproto.GetProperty(X, false, child, xproto.AtomWmClass, xproto.AtomString, 0, 1024).Reply()
		if err == nil && len(prop.Value) > 0 {
			if strings.Contains(strings.ToLower(string(prop.Value)), strings.ToLower(targetClass)) {
				match = true
			}
		}

		if match {
			attr, err := xproto.GetWindowAttributes(X, child).Reply()
			isViewable := (err == nil && attr.MapState == xproto.MapStateViewable)

			geom, err := xproto.GetGeometry(X, xproto.Drawable(child)).Reply()
			isLargeEnough := (err == nil && geom.Width > 10 && geom.Height > 10)

			if isViewable && isLargeEnough {
				return child
			}
		}
		if found := findWindow(X, child, targetClass); found != 0 {
			return found
		}
	}

	return 0
}

// x11ToImage converts the BGRA byte slice from X11 to an image.RGBA
func x11ToImage(data []byte, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// X11 usually sends BGRA or BGRA; standard Go image is RGBA.
	for i := 0; i < len(data); i += 4 {
		b := data[i]
		g := data[i+1]
		r := data[i+2]
		a := data[i+3]

		offset := i
		img.Pix[offset] = r
		img.Pix[offset+1] = g
		img.Pix[offset+2] = b
		img.Pix[offset+3] = a
	}
	return img
}
