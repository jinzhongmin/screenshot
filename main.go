package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"unsafe"

	"fmt"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/vova616/screenshot"
)

type Pointers struct {
	isInit bool
	startX int
	startY int
	endX   int
	endY   int
}
type App struct {
	mainWin     *gtk.Window
	shootButton *gtk.Button

	popWin   *gtk.Window
	popMenu  *gtk.Menu
	cancel   *gtk.MenuItem
	complete *gtk.MenuItem
	draw     *gtk.DrawingArea

	drawGc      *gdk.GC
	drawable    *gdk.Drawable
	shootPixbuf *gdkpixbuf.Pixbuf
}

func captureScreen() {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./shoot.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	f.Close()
}
func savePng(name string, x int, y int, xx int, yy int) {
	reader, err := os.Open("shoot.png")
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	sub := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(x, y, xx, yy))
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(file, sub)
}
func (p *Pointers) init() {
	p.startX = 0
	p.startY = 0
	p.endX = 0
	p.endY = 0
	p.isInit = true
}
func (p *Pointers) start(x int, y int) {
	p.startX = x
	p.startY = y
}
func (p *Pointers) end(x int, y int) {
	p.endX = x
	p.endY = y
}
func (p *Pointers) minWithDet() (int, int, int, int) {
	var x, y, dx, dy int
	if !p.isInit {
		fmt.Println("pointers no init!")
		return 0, 0, 0, 0
	}
	if p.startX < p.endX {
		x = p.startX
		dx = p.endX - p.startX
	} else {
		x = p.endX
		dx = p.startX - p.endX
	}
	if p.startY < p.endY {
		y = p.startY
		dy = p.endY - p.startY
	} else {
		y = p.endY
		dy = p.startY - p.endY
	}
	return x, y, dx, dy
}
func (p *Pointers) minWithMax() (int, int, int, int) {
	x, y, dx, dy := p.minWithDet()
	return x, y, x + dx, y + dy
}
func (a *App) Show() {
	a.mainWin.ShowAll()
}
func (a *App) creatMainWin() {
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("Screenshot")
	window.SetSizeRequest(200, 30)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})

	a.mainWin = window
}
func (a *App) creatUI() {
	a.shootButton = gtk.NewButtonWithLabel("截图")
	a.mainWin.Add(a.shootButton)

	var width, height int
	width = gdk.ScreenWidth()
	height = gdk.ScreenHeight()

	a.popWin = gtk.NewWindow(gtk.WINDOW_POPUP)
	a.popWin.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)
	a.popWin.SetSizeRequest(width, height)

	a.draw = gtk.NewDrawingArea()
	a.popWin.Add(a.draw)

	a.popMenu = gtk.NewMenu()
	a.complete = gtk.NewMenuItemWithLabel("完成")
	a.cancel = gtk.NewMenuItemWithLabel("取消")
	a.popMenu.Append(a.complete)
	a.popMenu.Append(a.cancel)

}
func (a *App) addEvent() {
	var width, height int
	width = gdk.ScreenWidth()
	height = gdk.ScreenHeight()

	press := false
	ps := Pointers{}
	ps.init()

	a.shootButton.Connect("clicked", func() {
		a.mainWin.Iconify()
		captureScreen()

		shootPixbuf, err := gdkpixbuf.NewPixbufFromFile("./shoot.png")
		if err != nil {
			log.Panic(err)
		}
		a.shootPixbuf = shootPixbuf
		a.popWin.ShowAll()
	})

	a.draw.Connect("expose-event", func() {
		a.drawable = a.draw.GetWindow().GetDrawable()
		a.drawGc = gdk.NewGC(a.drawable)
		a.drawable.DrawPixbuf(a.drawGc, a.shootPixbuf, 0, 0, 0, 0, width, height, gdk.RGB_DITHER_NONE, 0, 0)
	})

	a.draw.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		p := *(**gdk.EventButton)(unsafe.Pointer(&arg))
		if p.Button == 3 {
			a.popMenu.ShowAll()
			a.popMenu.Popup(nil, nil, func(menu *gtk.Menu, px, py *int, push_in *bool, data interface{}) {}, nil, 3, 1)
		} else {
			press = true
			ps.start(int(p.X), int(p.Y))
		}
	})
	a.draw.Connect("button-release-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		p := *(**gdk.EventButton)(unsafe.Pointer(&arg))
		if p.Button == 1 && press == true {
			press = false
			ps.end(int(p.X), int(p.Y))

			x, y, dx, dy := ps.minWithDet()

			a.drawable.DrawPixbuf(a.drawGc, a.shootPixbuf, 0, 0, 0, 0, width, height, gdk.RGB_DITHER_NONE, 0, 0)
			a.drawGc.SetRgbFgColor(gdk.NewColor("#003399"))

			a.drawable.DrawRectangle(a.drawGc, false, x, y, dx, dy)

			a.popMenu.ShowAll()
			a.popMenu.Popup(nil, nil, func(menu *gtk.Menu, px, py *int, push_in *bool, data interface{}) {}, nil, 3, 1)
		}
	})
	a.draw.Connect("motion-notify-event", func(ctx *glib.CallbackContext) {
		if press {
			arg := ctx.Args(0)
			p := *(**gdk.EventMotion)(unsafe.Pointer(&arg))

			ps.end(int(p.X), int(p.Y))
			x, y, dx, dy := ps.minWithDet()

			a.drawable.DrawPixbuf(a.drawGc, a.shootPixbuf, 0, 0, 0, 0, width, height, gdk.RGB_DITHER_NONE, 0, 0)
			a.drawGc.SetRgbFgColor(gdk.NewColor("#003399"))

			a.drawable.DrawRectangle(a.drawGc, false, x, y, dx, dy)
		}
	})
	a.cancel.Connect("activate", func() {
		a.popWin.Hide()
		a.mainWin.Deiconify()
	})
	a.complete.Connect("activate", func() {
		x, y, xx, yy := ps.minWithMax()
		a.popWin.Hide()
		filechooser := gtk.NewFileChooserDialog("save", a.popWin, gtk.FILE_CHOOSER_ACTION_SAVE, "保存", gtk.RESPONSE_ACCEPT)
		filechooser.Response(func() {
			savePng(filechooser.GetFilename(), x, y, xx, yy)
			filechooser.Destroy()
		})
		filechooser.Run()
	})
	a.draw.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK | gdk.BUTTON_RELEASE_MASK))
}
func main() {
	gtk.Init(&os.Args)

	app := App{}

	app.creatMainWin()
	app.creatUI()
	app.addEvent()

	app.Show()
	gtk.Main()
}
