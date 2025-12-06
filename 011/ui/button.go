package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BaseButton struct {
	parts.ControlBase      // 基本的なコントロール機能
	parts.MouseInteraction // マウス操作（クリックなど）の処理
	parts.Drawable         // 描画機能

	Color color.Color
}

// BaseButton生成
func NewBaseButton(x, y, w, h int) *BaseButton {
	b := &BaseButton{}
	b.InitBaseButton(nil, x, y, w, h)
	return b
}

func (b *BaseButton) InitBaseButton(g parts.AddChilder, x, y, w, h int) {
	b.InitControlBase(b, x, y, w, h)
	b.InitMouseInteraction(b)
	b.InitDrawable(b)
	if g != nil {
		g.AddChild(b)
	}
	b.OnDraw = b.drawBaseButton
	b.Color = color.RGBA{0x60, 0x60, 0x60, 0xff}
}

// BaseButtonの描画処理を行う
func (b *BaseButton) drawBaseButton(screen *ebiten.Image) {
	gx, gy := b.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), b.Color, false)
}

// Buttonはクリック可能なボタンコントロール
type Button struct {
	BaseButton         // マウス処理と描画機能を持つ基本ボタン
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
	b.InitBaseButton(nil, x, y, w, h)
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
	b.InitFocusable(b)
	if g != nil {
		g.AddChild(b)
	}

	b.OnPress = func() {
		b.Focus()
	}
}
