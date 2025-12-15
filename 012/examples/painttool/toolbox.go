package main

import (
	"MyProject/ui"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ToolButtonはツール選択ボタン
type ToolButton struct {
	ui.InteractiveControl

	Size int
}

// ToolButton初期化
func (b *ToolButton) InitToolButton(x, y, w, h, size int) {
	b.InitInteractiveControl(x, y, w, h)
	b.OnDraw = b.drawToolButton
	b.Size = size
}

// drawToolButtonはボタンの描画処理を行う
func (b *ToolButton) drawToolButton(screen *ebiten.Image) {
	gx, gy := b.GetGlobalPos()
	vector.StrokeRect(screen, float32(gx), float32(gy), float32(b.Width-1), float32(b.Height-1), 2, color.RGBA{0xff, 0xff, 0xff, 0xff}, false)
	vector.FillCircle(screen, float32(gx+b.Width/2), float32(gy+b.Height/2), float32(b.Size/2), color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
}

// ToolBoxはツール選択ボックス
type ToolBox struct {
	ui.GroupingControl

	marker   ui.Label
	Buttons  [8]ToolButton
	OnSelect func(int)
}

// ToolBox生成
func NewToolBox(x, y, w, h int) *ToolBox {
	tb := &ToolBox{}
	tb.InitGroupingControl(x, y, w, h)
	// ToolMarker初期化
	tb.marker.InitLabel(0, 2*32+8, 16, 16, "▶", 16)
	tb.AddChild(&tb.marker)

	// ToolButton初期化
	for i := range tb.Buttons {
		tb.Buttons[i].InitToolButton(18, i*32, 32, 32, i*2+2)
		tb.AddChild(&tb.Buttons[i])

		// ボタンクリック時の処理
		tb.Buttons[i].OnPress = func() {
			tb.marker.Y = i*32 + 8

			if tb.OnSelect != nil {
				tb.OnSelect(i*2 + 2)
			}
		}
	}

	return tb
}
