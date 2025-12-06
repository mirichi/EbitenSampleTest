package ui

import "MyProject/ui/parts"

type BlankWidget struct {
	parts.WidgetBase
	parts.Grouping
}

func NewBlankWidget(x, y, w, h int) *BlankWidget {
	c := &BlankWidget{}
	c.InitBlankWidget(nil, x, y, w, h)
	return c
}

func (c *BlankWidget) InitBlankWidget(g parts.AddChilder, x, y, w, h int) {
	c.InitWidgetBase(c, x, y, w, h)
	c.InitGrouping(c)
	if g != nil {
		g.AddChild(c)
	}
}
