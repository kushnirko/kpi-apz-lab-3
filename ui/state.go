package ui

import "image/color"

type Background struct {
	C color.Color
}

type BackgroundRectangle struct {
	X1, Y1, X2, Y2 float32
}

type Figure struct {
	X, Y float32
}

type StateData struct {
	Bg  *Background
	Br  *BackgroundRectangle
	Fgs []*Figure
}

func (sd *StateData) SetBgColor(c color.Color) {
	sd.Bg.C = c
}

func (sd *StateData) SetBr(br BackgroundRectangle) {
	sd.Br = &br
}

func (sd *StateData) AddFg(fg Figure) {
	sd.Fgs = append(sd.Fgs, &fg)
}

func (sd *StateData) ForEachFg(f func(*Figure)) {
	for _, fg := range sd.Fgs {
		f(fg)
	}
}

func (sd *StateData) Reset() {
	sd.SetBgColor(color.Black)
	sd.Br = nil
	sd.Fgs = nil
}

func (sd *StateData) GetBg() *Background {
	bg := *sd.Bg
	return &bg
}

func (sd *StateData) GetBr() *BackgroundRectangle {
	if sd.Br == nil {
		return nil
	}
	br := *sd.Br
	return &br
}

func (sd *StateData) GetFgs() []*Figure {
	fgs := make([]*Figure, len(sd.Fgs))
	for i, p := range sd.Fgs {
		fg := *p
		fgs[i] = &fg
	}
	return fgs
}

type Setter interface {
	SetBgColor(c color.Color)
	SetBr(br BackgroundRectangle)
	AddFg(fg Figure)
	ForEachFg(f func(*Figure))
	Reset()
}

type Getter interface {
	GetBg() *Background
	GetBr() *BackgroundRectangle
	GetFgs() []*Figure
}

type State interface {
	Setter
	Getter
}

func InitState() State {
	sd := StateData{
		Bg: &Background{C: color.Black},
	}
	return &sd
}
