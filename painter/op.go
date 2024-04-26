package painter

import (
	"image/color"

	"kpi-apz-lab-3/ui"
)

// Operation змінює вхідний стан вікна.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо стан вікна уже можна використовувати для відображення.
	Do(ui.StateSetter) bool
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(s ui.StateSetter) (ready bool) {
	for _, o := range ol {
		ready = o.Do(s) || ready
	}
	return
}

// UpdateOp операція, яка не змінює стану вікна, але сигналізує, що його потрібно розглядати як готового.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(ui.StateSetter) bool { return true }

// OperationFunc використовується для перетворення функції оновлення стану вікна в Operation.
type OperationFunc func(ui.StateSetter)

func (f OperationFunc) Do(s ui.StateSetter) bool {
	f(s)
	return false
}

// WhiteFill встановлює колір фону вікна у білий. Може бути використана як Operation через OperationFunc(WhiteFill).
func WhiteFill(s ui.StateSetter) {
	s.SetBgColor(color.White)
}

// GreenFill встановлює колір фону вікна у зелений. Може бути використана як Operation через OperationFunc(GreenFill).
func GreenFill(s ui.StateSetter) {
	s.SetBgColor(color.RGBA{G: 0xff, A: 0xff})
}

func BgRectDraw(s ui.StateSetter, x1, y1, x2, y2 float32) {
	s.SetBr(ui.BackgroundRectangle{X1: x1, Y1: y1, X2: x2, Y2: y2})
}

func FigureDraw(s ui.StateSetter, x, y float32) {
	s.AddFg(ui.Figure{X: x, Y: y})
}

func FiguresMove(s ui.StateSetter, x, y float32) {
	s.ForEachFg(func(fg *ui.Figure) {
		fg.X, fg.Y = x, y
	})
}

func UIStateReset(s ui.StateSetter) {
	s.Reset()
}
