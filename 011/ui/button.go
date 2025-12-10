package ui

import (
	"MyProject/ui/parts"
	"image/color"
)

// Buttonはクリック可能なボタン
type Button struct {
	InteractiveControl // マウス処理と描画機能を持つ基本セット
	parts.TextDrawable // テキスト描画機能
	parts.Focusable    // フォーカス制御機能

	HoverColor color.Color
	PressColor color.Color
}

// Button生成
func NewButton(x, y, w, h int, text string, size int) *Button {
	b := &Button{}
	b.InitButton(x, y, w, h, text, size)
	return b
}

func (b *Button) InitButton(x, y, w, h int, text string, size int) {
	b.InitInteractiveControl(x, y, w, h)
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
	b.InitFocusable(b)

	b.OnPress = func() {
		b.Focus()
	}

	b.HoverColor = color.RGBA{0x80, 0x80, 0x80, 0xff}
	b.PressColor = color.RGBA{0x40, 0x40, 0x40, 0xff}

	b.OnBeforeHandleInput = func() {
		b.BackColor = color.RGBA{0x60, 0x60, 0x60, 0xff}
	}

	b.OnHover = func() {
		b.BackColor = b.HoverColor
	}

	b.OnPressing = func() {
		b.BackColor = b.PressColor
	}
}
