package objects

import (
	"MyProject/gameparts"
	"MyProject/parts"
)

// Container is a generic entity that holds children using Grouping logic.
// It serves as a scene graph node.
type Container struct {
	parts.EntityBase
	parts.Grouping
	gameparts.Transform
}

func NewContainer(x, y float64) *Container {
	c := &Container{}
	c.InitContainer(x, y)
	return c
}

func (c *Container) InitContainer(x, y float64) {
	c.EntityBase.InitEntityBase(c, x, y)
	c.InitGrouping(c)
	c.InitTransform(c)
}

func (c *Container) GetGlobalPos() (float64, float64) {
	m := c.GetGlobalMatrix()
	return m.Apply(0, 0)
}
