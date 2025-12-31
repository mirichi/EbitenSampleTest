package ui

import (
	"MyProject/uiparts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Buttonはクリック可能なボタン
type Button struct {
	InteractiveControl   // マウス処理と描画機能を持つ基本セット
	uiparts.TextDrawable // テキスト描画機能
	uiparts.Focusable    // フォーカス制御機能

	BackColor color.Color
}

// Button生成
func NewButton(x, y, w, h float64, text string, size float64) *Button {
	b := &Button{}
	b.InitButton(x, y, w, h, text, size)
	return b
}

func (b *Button) InitButton(x, y, w, h float64, text string, size float64) {
	theme := uiparts.CurrentTheme

	b.InitInteractiveControl(x, y, w, h)
	b.InitTextDrawable(b, text, size, uiparts.AlignCenter, uiparts.AlignCenter, 0, 0, theme.Text, true)
	b.InitFocusable(b)

	b.OnPress = func() {
		b.Focus()
	}

	b.AddBeforeDrawFunction(func(screen *ebiten.Image) {
		gx, gy := b.GetGlobalPos()
		vector.FillRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), b.BackColor, false)
	})

	// Updateの最後で色を決定する
	b.AddAfterUpdateFunction(func() {
		theme := uiparts.CurrentTheme
		if b.IsHovering {
			b.BackColor = theme.ButtonHover
		} else if b.IsPressed {
			b.BackColor = theme.ButtonPress
		} else {
			b.BackColor = theme.ButtonNormal
		}
	})
}
