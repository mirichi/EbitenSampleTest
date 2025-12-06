package ui

import "MyProject/ui/parts"

type GroupingWidget struct {
	parts.WidgetBase
	parts.Grouping
}

func NewGroupingWidget(x, y, w, h int) *GroupingWidget {
	c := &GroupingWidget{}
	c.InitGroupingWidget(x, y, w, h)
	return c
}

func (c *GroupingWidget) InitGroupingWidget(x, y, w, h int) {
	c.InitWidgetBase(c, x, y, w, h)
	c.InitGrouping(c)
}
