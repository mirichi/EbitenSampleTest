package ui

import "MyProject/ui/parts"

type Blank struct {
	parts.ControlBase
	parts.Grouping
}

func NewBlank(x, y, w, h int) *Blank {
	c := &Blank{}
	c.InitControlBase(c, x, y, w, h)
	c.InitGrouping(c)
	return c
}
