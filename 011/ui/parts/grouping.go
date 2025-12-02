package parts

import (
	"MyProject/ui/input"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Groupingは子コントロールを管理する機能を提供します。
// 子コントロールの追加、入力イベントの伝播、更新、描画を統括します。
type Grouping struct {
	Control      Control
	Children     []Control
	OrderChange  bool
	AutoLayout   AutoLayoutInterface
	OnLayout     func()
	ClippingFlag bool
}

type internalLayout struct {
	grouping *Grouping
}

// LayoutはGroupingのデフォルトのレイアウト処理を行います。
// AutoResizableな子コントロールのサイズを親に合わせて調整します。
func (l *internalLayout) Layout() {
	for _, c := range l.grouping.Children {
		cb := c.GetControlBase()

		// ChildrenにAutoResizableがあった場合
		if cb.AutoResizable {
			// Groupingコントロールのサイズに合わせたサイズに更新する
			pcb := l.grouping.Control.GetControlBase()
			cb.Width = pcb.Width
			cb.Height = pcb.Height
		}

		// AutoResizableがGroupingだった場合、配下のレイアウト更新
		if al, ok := c.(AutoLayoutInterface); ok {
			al.Layout()
		}
	}
}

func (g *Grouping) InitGrouping(c Control) {
	g.Control = c
	g.AutoLayout = &internalLayout{g}
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
	c.GetControlBase().Parent = g.Control
	g.Children = append(g.Children, c)
}

// コントロールのHandleInput時に呼ばれるHandleInputFunction
func (c *Grouping) handleInputFunction(t input.Touch) bool {
	cb := c.Control.GetControlBase()
	gx, gy := cb.GetGlobalPos()
	x, y := t.Pos()
	if c.ClippingFlag {
		if gx <= x && x < gx+cb.Width && gy <= y && y < gy+cb.Height {
		} else {
			return false // 範囲に入っていなければ戻る
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
			return true
		}
	}

	return false
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (c *Grouping) updateFunction() {
	for i := len(c.Children) - 1; i >= 0; i-- {
		c.Children[i].Update()
	}
}

// コントロールのDraw時に呼ばれるDrawFunction
func (c *Grouping) drawFunction(screen *ebiten.Image) {
	cb := c.Control.GetControlBase()
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
	c.AutoLayout.Layout()
	if c.OnLayout != nil {
		c.OnLayout()
	}
}
