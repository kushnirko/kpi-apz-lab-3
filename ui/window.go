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
	Title         string
	Debug         bool
	OnScreenReady func()

	w    screen.Window
	st   chan State
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.st = make(chan State)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(st State) {
	pw.st <- st
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

	if pw.OnScreenReady != nil {
		pw.OnScreenReady()
	}

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

	var st State

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, &st)

		case st = <-pw.st:
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

func (pw *Visualizer) handleEvent(e any, st *State) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if st == nil {
			if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
				pw.drawDefaultUI(int(e.X), int(e.Y))
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if st == nil {
			centerX := pw.sz.Bounds().Dx() / 2
			centerY := pw.sz.Bounds().Dy() / 2
			pw.drawDefaultUI(centerX, centerY)
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.DrawUi(st)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI(centerX, centerY int) {
	pw.w.Fill(pw.sz.Bounds(), color.Black, draw.Src) // Фон.

	redColor := color.RGBA{R: 255, A: 255}

	horizontalRectWidth := 400
	horizontalRectHeight := 150
	verticalRectWidth := 130
	verticalRectHeight := 170

	horizontalRectX1 := centerX - horizontalRectWidth/2
	horizontalRectY1 := centerY - horizontalRectHeight
	horizontalRectX2 := centerX + horizontalRectWidth/2
	horizontalRectY2 := centerY
	horizontalRect := image.Rect(horizontalRectX1, horizontalRectY1, horizontalRectX2, horizontalRectY2)
	pw.w.Fill(horizontalRect, redColor, draw.Src)

	verticalRectX1 := centerX - verticalRectWidth/2
	verticalRectY1 := centerY
	verticalRectX2 := centerX + verticalRectWidth/2
	verticalRectY2 := centerY + verticalRectHeight
	verticalRect := image.Rect(verticalRectX1, verticalRectY1, verticalRectX2, verticalRectY2)
	pw.w.Fill(verticalRect, redColor, draw.Src)

	// Малювання білої рамки.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func (pw *Visualizer) TransformRelPoint(relX, relY float32) (int, int) { // rel - relative
	s := pw.sz.Size()
	x := int(relX * float32(s.X))
	y := int(relY * float32(s.Y))
	return x, y
}

func (pw *Visualizer) FillBg(c color.Color) {
	pw.w.Fill(pw.sz.Bounds(), c, draw.Src)
}

func (pw *Visualizer) DrawBgRect(x1, y1, x2, y2 int) {
	c := color.Black
	rect := image.Rect(x1, y1, x2, y2)
	pw.w.Fill(rect, c, draw.Src)
}

func (pw *Visualizer) DrawFigure(x, y int) {
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

func (pw *Visualizer) DrawUi(st *State) {
	pw.FillBg(st.Bg.C)

	if br := st.Br; br != nil {
		x1, y1 := pw.TransformRelPoint(br.X1, br.Y1)
		x2, y2 := pw.TransformRelPoint(br.X2, br.Y2)
		pw.DrawBgRect(x1, y1, x2, y2)
	}

	for _, f := range st.Fgs {
		x, y := pw.TransformRelPoint(f.X, f.Y)
		pw.DrawFigure(x, y)
	}
}
