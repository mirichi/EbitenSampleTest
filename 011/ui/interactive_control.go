package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// InteractiveControlはマウス操作と描画機能を持つ基本セット
type InteractiveControl struct {
	parts.ControlBase      // 基本機能
	parts.MouseInteraction // マウス操作（クリックなど）の処理
	parts.Drawable         // 描画機能

	BackColor color.Color
}

// InteractiveControl生成
func NewInteractiveControl(x, y, w, h int) *InteractiveControl {
	b := &InteractiveControl{}
	b.InitInteractiveControl(x, y, w, h)
	return b
}

// InteractiveControl初期化
func (b *InteractiveControl) InitInteractiveControl(x, y, w, h int) {
	b.InitControlBase(b, x, y, w, h)
	b.InitMouseInteraction(b)
	b.InitDrawable(b)

	b.OnDraw = func(screen *ebiten.Image) {
		gx, gy := b.GetGlobalPos()
		vector.FillRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), b.BackColor, false)
	}

	b.BackColor = color.RGBA{0x60, 0x60, 0x60, 0xff}
}
