package ui

import "MyProject/ui/parts"

type Blank struct {
	parts.ControlBase
	parts.Grouping
}

func NewBlank(x, y, w, h int) *Blank {
	c := &Blank{}
	c.InitBlank(nil, x, y, w, h)
	return c
}

func (c *Blank) InitBlank(g parts.GroupingInterface, x, y, w, h int) {
	c.InitControlBase(c, x, y, w, h)
	c.InitGrouping(c)
	if g != nil {
		g.AddChild(c)
	}
}
