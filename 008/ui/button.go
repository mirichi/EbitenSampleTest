package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Buttonはシンプルなボタンコントロール
type Button struct {
	*ControlBase
	*Clickable
	*Drawable
	*TextDrawable
}

// Button生成
// 現状ではグレーで固定
func NewButton(x, y, w, h int, text string, size int, f func()) *Button {
	image := ebiten.NewImage(w, h)
	image.Fill(color.RGBA{0x60, 0x60, 0x60, 0xff})

	cb := NewControlBase(x, y, w, h)
	b := &Button{
		ControlBase:  cb,
		Clickable:    NewClickable(cb, f),
		Drawable:     NewDrawable(cb, image),
		TextDrawable: NewTextDrawable(cb, text, size, AlignCenter, AlignCenter, 0, 0, color.White, true),
	}
	return b
}
