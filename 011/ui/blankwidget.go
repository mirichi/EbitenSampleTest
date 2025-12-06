package ui

import "MyProject/ui/parts"

type BlankWidget struct {
	parts.WidgetBase
	parts.Grouping
}

func NewBlankWidget(x, y, w, h int) *BlankWidget {
	c := &BlankWidget{}
	c.InitBlankWidget(x, y, w, h)
	return c
}

func (c *BlankWidget) InitBlankWidget(x, y, w, h int) {
	c.InitWidgetBase(c, x, y, w, h)
	c.InitGrouping(c)
}
