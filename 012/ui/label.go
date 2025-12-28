package ui

import (
	"MyProject/ui/parts"
	"image/color"
)

type Label struct {
	parts.ControlBase
	parts.TextDrawable
}

func NewLabel(x, y, w, h float64, text string, size float64) *Label {
	l := &Label{}
	l.InitLabel(x, y, w, h, text, size)
	return l
}

func (l *Label) InitLabel(x, y, w, h float64, text string, size float64) {
	l.InitControlBase(l, x, y, w, h)
	l.InitTextDrawable(l, text, size, parts.AlignLeft, parts.AlignCenter, 0, 0, color.White, true)
}
