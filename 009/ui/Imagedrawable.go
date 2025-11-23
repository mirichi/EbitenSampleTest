package ui

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
	ox, oy := cb.GetOwnerPos()
	op := &ebiten.DrawImageOptions{}
	// 親コントロールからの相対位置に描画
	op.GeoM.Translate(float64(ox+cb.X), float64(oy+cb.Y))
	screen.DrawImage(d.Image, op)
}
