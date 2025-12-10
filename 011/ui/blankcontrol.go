package ui

import (
	"MyProject/ui/parts"
)

type BlankControl struct {
	parts.ControlBase
}

func NewBlankControl(x, y, w, h int) *BlankControl {
	b := &BlankControl{}
	b.InitBlankControl(x, y, w, h)
	return b
}

func (b *BlankControl) InitBlankControl(x, y, w, h int) {
	b.InitControlBase(b, x, y, w, h)
}
