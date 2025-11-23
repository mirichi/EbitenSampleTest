package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Buttonはシンプルなボタンコントロール
type Button struct {
	*ControlBase
	*MouseInteraction
	*Drawable
	*TextDrawable
	*Focusable
}

// Button生成
// 現状ではグレーで固定
func NewButton(x, y, w, h int, text string, size int) *Button {
	b := &Button{}
	b.ControlBase = NewControlBase(b, x, y, w, h)
	b.MouseInteraction = NewMouseInteraction(b)
	b.Drawable = NewDrawable(b, b.drawButton)
	b.TextDrawable = NewTextDrawable(b, text, size, AlignCenter, AlignCenter, 0, 0, color.White, true)
	b.Focusable = NewFocusable(b)

	b.OnPress = func() {
		b.Focus()
	}

	return b
}

func (b *Button) drawButton(screen *ebiten.Image) {
	ox, oy := b.GetOwnerPos()
	vector.FillRect(screen, float32(ox+b.X), float32(oy+b.Y), float32(b.Width), float32(b.Height), color.RGBA{0x60, 0x60, 0x60, 0xff}, false)
	// Focusがある場合、枠を描画する
	if b.Focused {
		vector.StrokeRect(screen, float32(ox+b.X), float32(oy+b.Y), float32(b.Width), float32(b.Height), 2, color.RGBA{0xD0, 0xD0, 0xD0, 0xff}, false)
	}
}
