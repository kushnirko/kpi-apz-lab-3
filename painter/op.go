package painter

import (
	"github.com/roman-mazur/architecture-lab-3/ui"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує текстуру у білий колір. Може бути використана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує текстуру у зелений колір. Може бути використана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

func BgRectDraw(s ui.State, x1, y1, x2, y2 float32) {
	s.Br.X1, s.Br.Y1, s.Br.X2, s.Br.Y2 = x1, y1, x2, y2
}

func FigureDraw(s ui.State, x, y float32) {
	f := ui.Figure{X: x, Y: y}
	s.Fgs = append(s.Fgs, &f)
}

func FiguresMove(s ui.State, x, y float32) {
	for _, f := range s.Fgs {
		f.X, f.Y = x, y
	}
}

func TextureStateReset(s ui.State) {
	s.Bg.C = color.Black
	s.Br = nil
	s.Fgs = nil
}
