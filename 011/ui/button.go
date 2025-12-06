package ui

import (
	"MyProject/ui/parts"
	"image/color"
)

// Buttonはクリック可能なボタン
type Button struct {
	InteractiveWidget  // マウス処理と描画機能を持つ基本セット
	parts.TextDrawable // テキスト描画機能
	parts.Focusable    // フォーカス制御機能
}

// Button生成
func NewButton(x, y, w, h int, text string, size int) *Button {
	b := &Button{}
	b.InitButton(nil, x, y, w, h, text, size)
	return b
}

func (b *Button) InitButton(g parts.AddChilder, x, y, w, h int, text string, size int) {
	b.InitInteractiveWidget(nil, x, y, w, h)
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
	b.InitFocusable(b)

	if g != nil {
		g.AddChild(b)
	}

	b.OnPress = func() {
		b.Focus()
	}
}
