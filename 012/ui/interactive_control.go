package ui

import (
	"MyProject/uiparts"
)

// InteractiveControlはマウス操作と描画機能を持つ基本セット
type InteractiveControl struct {
	uiparts.ControlBase      // 基本機能
	uiparts.MouseInteraction // マウス操作（クリックなど）の処理

	// BackColor color.Color
}

// InteractiveControl生成
func NewInteractiveControl(x, y, w, h float64) *InteractiveControl {
	b := &InteractiveControl{}
	b.InitInteractiveControl(x, y, w, h)
	return b
}

// InteractiveControl初期化
func (b *InteractiveControl) InitInteractiveControl(x, y, w, h float64) {
	b.InitControlBase(b, x, y, w, h)
	b.InitMouseInteraction(b)

	// b.AddDrawFunction(func(screen *ebiten.Image) {
	// 	gx, gy := b.GetGlobalPos()
	// 	vector.FillRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), b.BackColor, false)
	// })

	// b.BackColor = color.RGBA{0x60, 0x60, 0x60, 0xff}
}
