package parts

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ImageDrawableはImageを描画する機能
type ImageDrawable struct {
	Control Control
	Image   *ebiten.Image
}

// ImageDrawable生成
func NewImageDrawable(c Control, image *ebiten.Image) *ImageDrawable {
	r := &ImageDrawable{
		Control: c,
		Image:   image,
	}

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetControlBase().AddDrawFunction(r.drawFunction)

	return r
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *ImageDrawable) drawFunction(screen *ebiten.Image) {
	cb := d.Control.GetControlBase()
	op := &ebiten.DrawImageOptions{}
	gx, gy := cb.GetGlobalPos()
	op.GeoM.Translate(float64(gx), float64(gy))
	screen.DrawImage(d.Image, op)
}
