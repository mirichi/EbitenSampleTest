package parts

import (
	"MyProject/ui/input"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Layouter interface {
	Layout()
}

type Sizer interface {
	GetSize() (float64, float64)
}

// Groupingは子コントロールを管理する機能
// 子コントロールの追加、入力イベントの伝播、更新、描画を統括する
type Grouping struct {
	Entity       Entity
	Children     []Entity
	OrderChange  bool
	AutoLayout   func(g *Grouping)
	OnLayout     func()
	ClippingFlag bool
}

func NewGrouping(e Entity) *Grouping {
	g := &Grouping{}
	g.InitGrouping(e)
	return g
}

func (g *Grouping) InitGrouping(e Entity) {
	g.Entity = e

	// コントロールのHandleInput時に呼ばれる関数を登録する
	g.Entity.GetEntityBase().AddHandleInputFunction(g.handleInputFunction)

	// コントロールのUpdate時に呼ばれる関数を登録する
	g.Entity.GetEntityBase().AddUpdateFunction(g.updateFunction)

	// コントロールのDraw時に呼ばれる関数を登録する
	g.Entity.GetEntityBase().AddDrawFunction(g.drawFunction)
}

// 子コントロールを登録する
func (g *Grouping) AddChild(c Entity) {
	c.GetEntityBase().Parent = g.Entity
	g.Children = append(g.Children, c)
}

// 子コントロールを削除する
func (g *Grouping) RemoveChild(c Entity) {
	for i, child := range g.Children {
		if child == c {
			g.Children = append(g.Children[:i], g.Children[i+1:]...)
			c.GetEntityBase().Parent = nil
			break
		}
	}
}

// コントロールのHandleInput時に呼ばれるHandleInputFunction
func (c *Grouping) handleInputFunction(t input.Touch) bool {
	handle := false
	eb := c.Entity.GetEntityBase()
	gx, gy := eb.GetGlobalPos()
	if t != nil {
		x, y := t.Pos()
		if c.ClippingFlag {
			// サイズを持っている（Sizerインターフェース実装）場合のみクリッピング判定を行う
			if sizer, ok := c.Entity.(Sizer); ok {
				w, h := sizer.GetSize()
				if gx <= x && x < gx+w && gy <= y && y < gy+h {
					// 範囲内
				} else {
					// 範囲に入っていなければ入力を渡さない
					t = nil
				}
			}
		}
	}

	if !eb.Visible {
		t = nil
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
	ox, oy := c.Entity.GetGlobalPos()
	if c.ClippingFlag {
		if sizer, ok := c.Entity.(Sizer); ok {
			w, h := sizer.GetSize()
			// クリッピング用SubImage。SubImageのSubImageは元の画像に対しての座標になるので入れ子構造でも大丈夫
			screen = screen.SubImage(image.Rect(int(ox), int(oy), int(ox+w), int(oy+h))).(*ebiten.Image)
		}
	}
	for _, ch := range c.Children {
		ch.Draw(screen)
	}
}

func (c *Grouping) Layout() {
	if c.AutoLayout != nil {
		c.AutoLayout(c)
		if c.OnLayout != nil {
			c.OnLayout()
		}
	} else {
		for _, ch := range c.Children {
			// Layouter実装時の再帰呼び出し
			if ali, ok := ch.(Layouter); ok {
				ali.Layout()
			}
		}

	}
}
