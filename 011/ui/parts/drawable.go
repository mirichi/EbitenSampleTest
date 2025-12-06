package parts

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Drawableは何かしらを描画する機能
type Drawable struct {
	Control Control
	OnDraw  func(screen *ebiten.Image)
}

func (d *Drawable) InitDrawable(c Control) {
	d.Control = c

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetControlBase().AddDrawFunction(d.drawFunction)
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *Drawable) drawFunction(screen *ebiten.Image) {
	if d.OnDraw != nil {
		d.OnDraw(screen)
	}
}
