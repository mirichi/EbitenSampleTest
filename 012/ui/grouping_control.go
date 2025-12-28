package ui

import "MyProject/ui/parts"

// GroupingControlは複合コントロールの親になるコントロール
type GroupingControl struct {
	parts.ControlBase
	parts.Grouping
}

func NewGroupingControl(x, y, w, h float64) *GroupingControl {
	c := &GroupingControl{}
	c.InitGroupingControl(x, y, w, h)
	return c
}

func (c *GroupingControl) InitGroupingControl(x, y, w, h float64) {
	c.InitControlBase(c, x, y, w, h)
	c.InitGrouping(c)
}
