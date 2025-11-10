package ui

import "github.com/hajimehoshi/ebiten/v2"

// Groupingは子コントロールを保持する機能
type Grouping struct {
	control  *ControlBase
	Childlen []Control
}

// Grouping生成
func NewGrouping(c *ControlBase) *Grouping {
	r := &Grouping{control: c}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.AddUpdateFunction(r.UpdateFunction)

	// コントロールのDraw時に呼ばれる関数を登録する
	c.AddDrawFunction(r.DrawFunction)

	return r
}

// 子コントロールを登録する
func (g *Grouping) AddChild(c Control) {
	c.GetControlBase().Parent = g.control
	g.Childlen = append(g.Childlen, c)
}

// コントロールのUpdate時に呼ばれるUpdateFunction
// 子コントロールのUpdateを呼ぶ
func (c *Grouping) UpdateFunction(t TouchInfo) {
	for _, ch := range c.Childlen {
		ch.Update(t)
	}
}

// コントロールのDraw時に呼ばれるDrawFunction
// 子コントロールのDrawを呼ぶ
func (c *Grouping) DrawFunction(screen *ebiten.Image) {
	for _, ch := range c.Childlen {
		ch.Draw(screen)
	}
}
