package parts

import (
	"MyProject/ui/input"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Layouter interface {
	Layout()
}

// Groupingは子コントロールを管理する機能を提供します。
// 子コントロールの追加、入力イベントの伝播、更新、描画を統括します。
type Grouping struct {
	Widget       Control
	Children     []Control
	OrderChange  bool
	AutoLayout   func(g *Grouping)
	OnLayout     func()
	ClippingFlag bool
}

func NewGrouping(c Control) *Grouping {
	g := &Grouping{}
	g.InitGrouping(c)
	return g
}

func (g *Grouping) InitGrouping(c Control) {
	g.Widget = c
	g.AutoLayout = DefaultLayout
	g.ClippingFlag = true

	// コントロールのHandleInput時に呼ばれる関数を登録する
	c.GetControlBase().AddHandleInputFunction(g.handleInputFunction)

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(g.updateFunction)

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetControlBase().AddDrawFunction(g.drawFunction)
}

// 子コントロールを登録する
func (g *Grouping) AddChild(c Control) {
	c.GetControlBase().Parent = g.Widget
	g.Children = append(g.Children, c)
}

// コントロールのHandleInput時に呼ばれるHandleInputFunction
func (c *Grouping) handleInputFunction(t input.Touch) bool {
	handle := false
	cb := c.Widget.GetControlBase()
	gx, gy := cb.GetGlobalPos()
	if t != nil {
		x, y := t.Pos()
		if c.ClippingFlag {
			if gx <= x && x < gx+cb.Width && gy <= y && y < gy+cb.Height {
			} else {
				// 範囲に入っていなければ入力を渡さない
				t = nil
			}
		}
	}

	for i := len(c.Children) - 1; i >= 0; i-- {
		if c.Children[i].HandleInput(t) {
			// 順序変更は押された時のみ行う（ホバーで手前に来ないように）
			if c.OrderChange && t.IsJustPressed() {
				// スライスからiを消して最後に追加する
				c.Children = append(c.Children, c.Children[i])
				c.Children = append(c.Children[:i], c.Children[i+1:]...)
			}

			// 入力を処理したら消費する
			t = nil
			handle = true
		}
	}

	return handle
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (c *Grouping) updateFunction() {
	for i := len(c.Children) - 1; i >= 0; i-- {
		c.Children[i].Update()
	}
}

// コントロールのDraw時に呼ばれるDrawFunction
func (c *Grouping) drawFunction(screen *ebiten.Image) {
	cb := c.Widget.GetControlBase()
	ox, oy := cb.GetGlobalPos()
	if c.ClippingFlag {
		// クリッピング用SubImage。SubImageのSubImageは元の画像に対しての座標になるので入れ子構造でも大丈夫
		screen = screen.SubImage(image.Rect(ox, oy, ox+cb.Width, oy+cb.Height)).(*ebiten.Image)
	}
	for _, ch := range c.Children {
		ch.Draw(screen)
	}
}

func (c *Grouping) Layout() {
	c.AutoLayout(c)
	if c.OnLayout != nil {
		c.OnLayout()
	}
}
