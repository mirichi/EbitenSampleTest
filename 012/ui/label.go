package ui

import (
	"MyProject/ui/parts"
	"image/color"
)

type Label struct {
	parts.ControlBase
	parts.TextDrawable
}

func NewLabel(x, y, w, h int, text string, size int) *Label {
	l := &Label{}
	l.InitLabel(x, y, w, h, text, size)
	return l
}

func (l *Label) InitLabel(x, y, w, h int, text string, size int) {
	l.InitControlBase(l, x, y, w, h)
	l.InitTextDrawable(l, text, size, parts.AlignLeft, parts.AlignCenter, 0, 0, color.White, true)
}
