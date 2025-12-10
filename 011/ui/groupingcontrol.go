package ui

import "MyProject/ui/parts"

type GroupingControl struct {
	parts.ControlBase
	parts.Grouping
}

func NewGroupingControl(x, y, w, h int) *GroupingControl {
	c := &GroupingControl{}
	c.InitGroupingControl(x, y, w, h)
	return c
}

func (c *GroupingControl) InitGroupingControl(x, y, w, h int) {
	c.InitControlBase(c, x, y, w, h)
	c.InitGrouping(c)
}
