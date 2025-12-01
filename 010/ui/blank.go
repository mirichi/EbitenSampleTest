package ui

import "MyProject/ui/parts"

type Blank struct {
	*parts.ControlBase
	*parts.Grouping
}

func NewBlank(x, y, w, h int) *Blank {
	c := &Blank{}
	c.ControlBase = parts.NewControlBase(c, x, y, w, h)
	c.Grouping = parts.NewGrouping(c)
	return c
}
