package parts

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ImageDrawableはImageを描画する機能
type ImageDrawable struct {
	Widget Widget
	Image  *ebiten.Image
}

func (d *ImageDrawable) InitImageDrawable(c Widget, image *ebiten.Image) {
	d.Widget = c
	d.Image = image

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetWidgetBase().AddDrawFunction(d.drawFunction)
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *ImageDrawable) drawFunction(screen *ebiten.Image) {
	cb := d.Widget.GetWidgetBase()
	op := &ebiten.DrawImageOptions{}
	gx, gy := cb.GetGlobalPos()
	op.GeoM.Translate(float64(gx), float64(gy))
	screen.DrawImage(d.Image, op)
}
