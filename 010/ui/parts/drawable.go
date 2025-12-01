package parts

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Drawableは何かしらを描画する機能
type Drawable struct {
	Control Control
	OnDraw  func(screen *ebiten.Image)
}

// Drawable生成
func NewDrawable(c Control, f func(screen *ebiten.Image)) *Drawable {
	r := &Drawable{
		Control: c,
		OnDraw:  f,
	}

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetControlBase().AddDrawFunction(r.drawFunction)

	return r
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *Drawable) drawFunction(screen *ebiten.Image) {
	d.OnDraw(screen)
}
