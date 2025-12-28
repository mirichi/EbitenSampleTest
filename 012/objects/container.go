package objects

import (
	"MyProject/parts"

	"github.com/hajimehoshi/ebiten/v2"
)

// Container is a generic entity that holds children using Grouping logic.
// It serves as a scene graph node.
type Container struct {
	parts.EntityBase
	parts.Grouping
	Angle  float64
	ScaleX float64
	ScaleY float64
}

func NewContainer(x, y float64) *Container {
	c := &Container{}
	c.InitContainer(x, y)
	return c
}

func (c *Container) InitContainer(x, y float64) {
	c.EntityBase.InitEntityBase(c, x, y)
	c.InitGrouping(c)
	c.ScaleX = 1
	c.ScaleY = 1
}

func (c *Container) GetChildren() []parts.Entity {
	return c.Children
}

// GlobalMatrixer allows checking if a parent supports GetGlobalMatrix
type GlobalMatrixer interface {
	GetGlobalMatrix() ebiten.GeoM
}

func (c *Container) GetGlobalMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	// Local Transform: Scale -> Rotate -> Translate
	m.Scale(c.ScaleX, c.ScaleY)
	m.Rotate(c.Angle)
	m.Translate(c.X, c.Y)

	// Parent Transform
	if p := c.EntityBase.Parent; p != nil {
		if gm, ok := p.(GlobalMatrixer); ok {
			// Parent has global matrix support
			pm := gm.GetGlobalMatrix()
			m.Concat(pm)
		} else {
			// Fallback: Parent only has Position
			px, py := p.GetGlobalPos()
			m.Translate(px, py)
		}
	}

	return m
}

func (c *Container) GetGlobalPos() (float64, float64) {
	m := c.GetGlobalMatrix()
	// Extract translation from matrix
	// Element 2,0 and 2,1 are Tx, Ty in Ebiten's GeoM (technically internal implementation)
	// Safe way: Apply to (0,0)
	return m.Apply(0, 0)
}
