package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Buttonはクリック可能なボタンコントロール
type Button struct {
	parts.ControlBase      // 基本的なコントロール機能
	parts.MouseInteraction // マウス操作（クリックなど）の処理
	parts.Drawable         // 描画機能
	parts.TextDrawable     // テキスト描画機能
	parts.Focusable        // フォーカス制御機能

	Color color.Color
}

// Button生成
func NewButton(x, y, w, h int, text string, size int) *Button {
	b := &Button{}
	b.InitControlBase(b, x, y, w, h)
	b.InitMouseInteraction(b)
	b.InitDrawable(b, b.drawButton)
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
	b.InitFocusable(b)
	b.Color = color.RGBA{0x60, 0x60, 0x60, 0xff}

	b.OnPress = func() {
		b.Focus()
	}

	return b
}

// drawButtonはボタンの描画処理を行う
func (b *Button) drawButton(screen *ebiten.Image) {
	gx, gy := b.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), b.Color, false)
}
