package ui

import (
	"MyProject/ui/parts"
)

// Buttonはクリック可能なボタン
type Button struct {
	InteractiveControl // マウス処理と描画機能を持つ基本セット
	parts.TextDrawable // テキスト描画機能
	parts.Focusable    // フォーカス制御機能
}

// Button生成
func NewButton(x, y, w, h int, text string, size int) *Button {
	b := &Button{}
	b.InitButton(x, y, w, h, text, size)
	return b
}

func (b *Button) InitButton(x, y, w, h int, text string, size int) {
	theme := parts.CurrentTheme

	b.InitInteractiveControl(x, y, w, h)
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, theme.Text, true)
	b.InitFocusable(b)

	b.OnPress = func() {
		b.Focus()
	}

	// Updateの最後で色を決定する
	b.AddAfterUpdateFunction(func() {
		theme := parts.CurrentTheme
		if b.IsHovering {
			b.BackColor = theme.ButtonHover
		} else if b.IsPressed {
			b.BackColor = theme.ButtonPress
		} else {
			b.BackColor = theme.ButtonNormal
		}
	})
}
