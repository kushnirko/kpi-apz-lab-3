package painter

import (
	"image/color"
	"reflect"
	"testing"

	"kpi-apz-lab-3/ui"
)

func TestLoop_Post(t *testing.T) {
	var (
		l       Loop
		tr      testReceiver
		testOps []string
	)
	l.Receiver = &tr

	l.Start()

	// Перевіряємо чи коректно встановиться колір фону
	l.Post(OperationFunc(GreenFill))
	l.Post(OperationFunc(WhiteFill))

	// Перевіряємо асинхронність виконання операцій
	l.Post(OperationFunc(func(ui.StateSetter) {
		testOps = append(testOps, "op 1")
		l.Post(OperationFunc(func(ui.StateSetter) {
			testOps = append(testOps, "op 4")
		}))
	}))
	l.Post(OperationFunc(func(ui.StateSetter) {
		testOps = append(testOps, "op 2")
	}))

	// Перевіряємо чи коректно додаються операції з інших рутин
	go func() {
		l.Post(OperationFunc(func(ui.StateSetter) {
			testOps = append(testOps, "op 3")
			l.Post(OperationFunc(func(ui.StateSetter) {
				testOps = append(testOps, "op 5")
			}))
		}))
		l.Post(OperationFunc(func(s ui.StateSetter) {
			BgRectDraw(s, 0.1, 0.2, 0.3, 0.4)
		}))
	}()

	l.Post(UpdateOp)

	l.StopAndWait()

	if tr.lastState == nil {
		t.Fatal("window state was not updated")
	}

	if c := tr.lastState.GetBg().C; c != color.White {
		t.Error("background color is not white:", c)
	}

	if !reflect.DeepEqual(testOps, []string{"op 1", "op 2", "op 3", "op 4", "op 5"}) {
		t.Error("bad operations order:", testOps)
	}

	resBr := tr.lastState.GetBr()
	expectedBr := ui.BackgroundRectangle{X1: 0.1, Y1: 0.2, X2: 0.3, Y2: 0.4}
	if !reflect.DeepEqual(*resBr, expectedBr) {
		t.Error("bad background rectangle:", resBr)
	}

	if !l.mq.empty() {
		t.Error("message query not empty:", l.mq.ops)
	}
}

type testReceiver struct {
	lastState ui.StateGetter
}

func (tr *testReceiver) Update(g ui.StateGetter) {
	tr.lastState = g
}
