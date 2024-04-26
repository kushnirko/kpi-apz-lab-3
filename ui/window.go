package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title string
	Debug bool

	w    screen.Window
	g    chan StateGetter
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.g = make(chan StateGetter)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(g StateGetter) {
	pw.g <- g
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Width:  800,
		Height: 800,
		Title:  pw.Title,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var g StateGetter

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, g)

		case g = <-pw.g:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, g StateGetter) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if g == nil {
			if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
				pw.drawDefaultUI(int(e.X), int(e.Y))
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if g == nil {
			centerX := pw.sz.Bounds().Dx() / 2
			centerY := pw.sz.Bounds().Dy() / 2
			pw.drawDefaultUI(centerX, centerY)
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.drawUI(g)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI(centerX, centerY int) {
	pw.fillBg(color.Black)
	pw.drawFigure(centerX, centerY)

	// Малювання білої рамки.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func (pw *Visualizer) transformRelPoint(relX, relY float32) (int, int) { // rel - relative
	s := pw.sz.Size()
	x := int(relX * float32(s.X))
	y := int(relY * float32(s.Y))
	return x, y
}

func (pw *Visualizer) fillBg(c color.Color) {
	pw.w.Fill(pw.sz.Bounds(), c, draw.Src)
}

func (pw *Visualizer) drawBgRect(x1, y1, x2, y2 int) {
	c := color.Black
	rect := image.Rect(x1, y1, x2, y2)
	pw.w.Fill(rect, c, draw.Src)
}

func (pw *Visualizer) drawFigure(x, y int) {
	c := color.RGBA{R: 0xff, A: 0xff}

	// h - horizontal, v - vertical
	// W - width, H - height
	hRectW := 400
	hRectH := 150
	vRectW := 130
	vRectH := 170

	hRect := image.Rect(x-hRectW/2, y-hRectH, x+hRectW/2, y)
	pw.w.Fill(hRect, c, draw.Src)

	vRect := image.Rect(x-vRectW/2, y, x+vRectW/2, y+vRectH)
	pw.w.Fill(vRect, c, draw.Src)
}

func (pw *Visualizer) drawUI(g StateGetter) {
	pw.fillBg(g.GetBg().C)

	if br := g.GetBr(); br != nil {
		x1, y1 := pw.transformRelPoint(br.X1, br.Y1)
		x2, y2 := pw.transformRelPoint(br.X2, br.Y2)
		pw.drawBgRect(x1, y1, x2, y2)
	}

	for _, f := range g.GetFgs() {
		x, y := pw.transformRelPoint(f.X, f.Y)
		pw.drawFigure(x, y)
	}
}
