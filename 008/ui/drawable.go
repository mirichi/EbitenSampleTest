package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// DrawableはImageを描画する機能
type Drawable struct {
	control *ControlBase
	Image   *ebiten.Image
}

// Drawable生成
func NewDrawable(c *ControlBase, image *ebiten.Image) *Drawable {
	r := &Drawable{
		control: c,
		Image:   image,
	}

	// コントロールのDraw時に呼ばれる関数を登録する
	c.AddDrawFunction(r.DrawFunction)

	return r
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *Drawable) DrawFunction(screen *ebiten.Image) {
	ox, oy := d.control.GetOwnerPos()
	op := &ebiten.DrawImageOptions{}
	// 親コントロールからの相対位置に描画
	op.GeoM.Translate(float64(ox+d.control.X), float64(oy+d.control.Y))
	screen.DrawImage(d.Image, op)
}
