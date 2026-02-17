package internal

import (
	"image"
	"log"
	"strings"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func screenshot() image.Image {
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
