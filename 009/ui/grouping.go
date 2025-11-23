package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Groupingは子コントロールを保持する機能
type Grouping struct {
	Control     Control
	Children    []Control
	OrderChange bool
	AutoLayout  AutoLayoutInterface
}

type internalLayout struct {
	grouping *Grouping
}

// GroupingがLayoutメソッドに反応するためのデフォルトの仕掛け
func (l *internalLayout) Layout() {
	for _, c := range l.grouping.Children {
		cb := c.GetControlBase()

		// ChildrenにAutoResizableがあった場合
		if cb.AutoResizable {
			// Groupingコントロールのサイズに合わせたサイズに更新する
			pcb := l.grouping.Control.GetControlBase()
			cb.Width = pcb.Width
			cb.Height = pcb.Height

			// 配下のレイアウト更新
			l.grouping.Layout()
		}
	}
}

// Grouping生成
func NewGrouping(c Control) *Grouping {
	r := &Grouping{Control: c}
	r.AutoLayout = &internalLayout{r}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(r.updateFunction)

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetControlBase().AddDrawFunction(r.drawFunction)

	return r
}

// 子コントロールを登録する
func (g *Grouping) AddChild(c Control) {
	c.GetControlBase().Parent = g.Control
	g.Children = append(g.Children, c)
}

// コントロールのUpdate時に呼ばれるUpdateFunction
// 子コントロールのUpdateを呼ぶ
func (c *Grouping) updateFunction(t TouchInfo) UpdateResult {
	var r UpdateResult = NotConsumed
	j := -1

	// 押されたかどうかを判定しておく（ループ内でtがnilになる可能性があるため）
	isJustPressed := t != nil && t.IsJustPressed()

	for i := len(c.Children) - 1; i >= 0; i-- {
		ch := c.Children[i]
		if ch.Update(t) == Consumed {
			// Groupingが消費した場合、コントロール内の伝搬も停止する
			r = StopPropagation
			t = nil
			j = i
		}
	}

	if j >= 0 {
		// 順序変更は押された時のみ行う（ホバーで手前に来ないように）
		if c.OrderChange && isJustPressed {
			// スライスからjを消して最後に追加する
			c.Children = append(c.Children, c.Children[j])
			c.Children = append(c.Children[:j], c.Children[j+1:]...)
		}
	}

	return r
}

// コントロールのDraw時に呼ばれるDrawFunction
// 子コントロールのDrawを呼ぶ
func (c *Grouping) drawFunction(screen *ebiten.Image) {
	// cb := c.Control.GetControlBase()
	// ox, oy := cb.GetOwnerPos()
	// for _, ch := range c.Children {
	// 	// クリッピング用SubImage。SubImageのSubImageは元の画像に対しての座標になるので入れ子構造でも大丈夫
	// 	subimage := screen.SubImage(image.Rect(cb.X+ox, cb.Y+oy, cb.X+ox+cb.Width, cb.Y+oy+cb.Height)).(*ebiten.Image)
	// 	ch.Draw(subimage)
	// }

	// ブラウザでは何も描画できなくなってしまうのでこっちを使う
	for _, ch := range c.Children {
		ch.Draw(screen)
	}
}

func (c *Grouping) Layout() {
	c.AutoLayout.Layout()
}
