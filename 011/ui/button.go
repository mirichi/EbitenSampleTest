package ui

import (
	"MyProject/ui/parts"
	"image/color"
)

// Buttonはクリック可能なボタン
type Button struct {
	InteractiveControl // マウス処理と描画機能を持つ基本セット
	parts.Updatable
	parts.TextDrawable // テキスト描画機能
	parts.Focusable    // フォーカス制御機能

	NormalColor color.Color
	HoverColor  color.Color
	PressColor  color.Color
}

// Button生成
func NewButton(x, y, w, h int, text string, size int) *Button {
	b := &Button{}
	b.InitButton(x, y, w, h, text, size)
	return b
}

func (b *Button) InitButton(x, y, w, h int, text string, size int) {
	b.InitInteractiveControl(x, y, w, h)
	b.InitUpdatable(b)
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
	b.InitFocusable(b)

	b.OnPress = func() {
		b.Focus()
	}

	b.NormalColor = color.RGBA{0x60, 0x60, 0x60, 0xff}
	b.HoverColor = color.RGBA{0x80, 0x80, 0x80, 0xff}
	b.PressColor = color.RGBA{0x40, 0x40, 0x40, 0xff}

	b.OnUpdate = func() {
		if b.IsHovering {
			b.BackColor = b.HoverColor
		} else if b.IsPressed {
			b.BackColor = b.PressColor
		} else {
			b.BackColor = b.NormalColor
		}
	}
}
