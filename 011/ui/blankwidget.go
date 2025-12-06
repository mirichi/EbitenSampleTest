package ui

import (
	"MyProject/ui/parts"
)

type BlankWidget struct {
	parts.WidgetBase
}

func NewBlankWidget(x, y, w, h int) *BlankWidget {
	b := &BlankWidget{}
	b.InitBlankWidget(x, y, w, h)
	return b
}

func (b *BlankWidget) InitBlankWidget(x, y, w, h int) {
	b.InitWidgetBase(b, x, y, w, h)
}
