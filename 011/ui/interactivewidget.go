package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// InteractiveWidgetはマウス操作と描画機能を持つ基本セット
type InteractiveWidget struct {
	parts.WidgetBase       // 基本機能
	parts.MouseInteraction // マウス操作（クリックなど）の処理
	parts.Drawable         // 描画機能

	Color color.Color
}

// InteractiveWidget生成
func NewInteractiveWidget(x, y, w, h int) *InteractiveWidget {
	b := &InteractiveWidget{}
	b.InitInteractiveWidget(nil, x, y, w, h)
	return b
}

// InteractiveWidget初期化
func (b *InteractiveWidget) InitInteractiveWidget(g parts.AddChilder, x, y, w, h int) {
	b.InitWidgetBase(b, x, y, w, h)
	b.InitMouseInteraction(b)
	b.InitDrawable(b)
	if g != nil {
		g.AddChild(b)
	}
	b.OnDraw = b.drawInteractiveWidget
	b.Color = color.RGBA{0x60, 0x60, 0x60, 0xff}
}

// InteractiveWidgetの描画処理を行う
func (b *InteractiveWidget) drawInteractiveWidget(screen *ebiten.Image) {
	gx, gy := b.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), b.Color, false)
}
