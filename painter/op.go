package painter

import (
	"github.com/roman-mazur/architecture-lab-3/ui"
	"image/color"
)

// Operation змінює вхідний стан вікна.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо стан вікна уже можна використовувати для відображення.
	Do(s ui.State) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(s ui.State) (ready bool) {
	for _, o := range ol {
		ready = o.Do(s) || ready
	}
	return
}

// UpdateOp операція, яка не змінює стану вікна, але сигналізує, що його потрібно розглядати як готового.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(ui.State) bool { return true }

// OperationFunc використовується для перетворення функції оновлення стану вікна в Operation.
type OperationFunc func(ui.State)

func (f OperationFunc) Do(s ui.State) bool {
	f(s)
	return false
}

// WhiteFill встановлює колір фону вікна у білий. Може бути використана як Operation через OperationFunc(WhiteFill).
func WhiteFill(s ui.State) {
	s.Bg.C = color.White
}

// GreenFill встановлює колір фону вікна у зелений. Може бути використана як Operation через OperationFunc(GreenFill).
func GreenFill(s ui.State) {
	s.Bg.C = color.RGBA{G: 0xff, A: 0xff}
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
