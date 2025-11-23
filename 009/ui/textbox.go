package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TextBox struct {
	*ControlBase
	*MouseInteraction
	*Drawable
	*TextInputable
	*Focusable

	field   textinput.Field
	counter int
}

func NewTextBox(x, y, w, h int, text string, fontSize int) *TextBox {
	tb := &TextBox{}

	tb.ControlBase = NewControlBase(tb, x, y, w, h)

	// 初期テキスト
	tb.field.SetTextAndSelection(text, 0, len(text))

	tb.MouseInteraction = NewMouseInteraction(tb)
	tb.OnPress = func() {
		tb.Focus()
	}

	tb.Drawable = NewDrawable(tb, tb.drawTextBox)
	tb.TextInputable = NewTextInputable(tb, &tb.field, &tb.counter, fontSize, AlignLeft, color.White, color.White, func() bool { return tb.Focused })
	tb.Focusable = NewFocusable(tb)
	tb.OnFocus = func() {
		tb.field.Focus()
	}
	tb.OnBlur = func() {
		tb.field.Blur()
	}

	return tb
}

func (tb *TextBox) drawTextBox(screen *ebiten.Image) {
	ox, oy := tb.GetOwnerPos()
	x := float32(ox + tb.X)
	y := float32(oy + tb.Y)

	// 背景と枠線（自前描画）
	w := float32(tb.Width)
	h := float32(tb.Height)
	vector.FillRect(screen, x, y, w, h, color.RGBA{0x60, 0x60, 0x60, 0xff}, false)
	vector.StrokeRect(screen, x, y, w, h, 2, color.RGBA{0x80, 0x80, 0x80, 0xff}, false)

	if tb.Focused {
		vector.StrokeRect(screen, x, y, w, h, 2, color.RGBA{0xD0, 0xD0, 0xD0, 0xff}, false)
	}
}
