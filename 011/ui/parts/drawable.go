package parts

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Drawableは何かしらを描画する機能
type Drawable struct {
	Widget Widget
	OnDraw func(screen *ebiten.Image)
}

func (d *Drawable) InitDrawable(c Widget) {
	d.Widget = c

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetWidgetBase().AddDrawFunction(d.drawFunction)
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *Drawable) drawFunction(screen *ebiten.Image) {
	if d.OnDraw != nil {
		d.OnDraw(screen)
	}
}
