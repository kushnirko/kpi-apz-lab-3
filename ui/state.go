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

type State struct {
	Bg  *Background
	Br  *BackgroundRectangle
	Fgs []*Figure
}
