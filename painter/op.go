package painter

import (
	"image/color"
	"kpi-apz-lab-3/ui"
)

// Operation змінює вхідний стан вікна.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо стан вікна уже можна використовувати для відображення.
	Do(st ui.State) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(st ui.State) (ready bool) {
	for _, o := range ol {
		ready = o.Do(st) || ready
	}
	return
}

// UpdateOp операція, яка не змінює стану вікна, але сигналізує, що його потрібно розглядати як готового.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(ui.State) bool { return true }

// OperationFunc використовується для перетворення функції оновлення стану вікна в Operation.
type OperationFunc func(ui.State)

func (f OperationFunc) Do(st ui.State) bool {
	f(st)
	return false
}

// WhiteFill встановлює колір фону вікна у білий. Може бути використана як Operation через OperationFunc(WhiteFill).
func WhiteFill(st ui.State) {
	st.Bg.C = color.White
}

// GreenFill встановлює колір фону вікна у зелений. Може бути використана як Operation через OperationFunc(GreenFill).
func GreenFill(st ui.State) {
	st.Bg.C = color.RGBA{G: 0xff, A: 0xff}
}

func BgRectDraw(st ui.State, x1, y1, x2, y2 float32) {
	st.Br.X1, st.Br.Y1, st.Br.X2, st.Br.Y2 = x1, y1, x2, y2
}

func FigureDraw(st ui.State, x, y float32) {
	f := ui.Figure{X: x, Y: y}
	st.Fgs = append(st.Fgs, &f)
}

func FiguresMove(st ui.State, x, y float32) {
	for _, f := range st.Fgs {
		f.X, f.Y = x, y
	}
}

func UIStateReset(st ui.State) {
	st.Bg.C = color.Black
	st.Br = nil
	st.Fgs = nil
}
