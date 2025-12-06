package main

import (
	"MyProject/ui"
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ToolButtonはツール選択ボタン
type ToolButton struct {
	ui.InteractiveWidget

	Size int
}

// ToolButton初期化
func (b *ToolButton) InitToolButton(x, y, w, h, size int) {
	b.InitInteractiveWidget(x, y, w, h)
	b.OnDraw = b.drawToolButton
	b.Size = size
}

// drawToolButtonはボタンの描画処理を行う
func (b *ToolButton) drawToolButton(screen *ebiten.Image) {
	gx, gy := b.GetGlobalPos()
	vector.StrokeRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), 2, color.RGBA{0xff, 0xff, 0xff, 0xff}, false)
	vector.FillCircle(screen, float32(gx+b.Width/2), float32(gy+b.Height/2), float32(b.Size/2), color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
}

// ToolBoxはツール選択ボックス
type ToolBox struct {
	parts.WidgetBase
	parts.Grouping

	marker   ui.Label
	buttons  [8]ToolButton
	OnSelect func(int)
}

// ToolBox生成
func NewToolBox(x, y, w, h int) *ToolBox {
	tb := &ToolBox{}
	tb.InitWidgetBase(tb, x, y, w, h)
	tb.InitGrouping(tb)
	// ToolMarker初期化
	tb.marker.InitLabel(0, 2*32+8, 16, 16, "▶", 16)
	tb.AddChild(&tb.marker)

	// ToolButton初期化
	for i := range tb.buttons {
		tb.buttons[i].InitToolButton(18, i*32, 32, 32, i*2+2)
		tb.AddChild(&tb.buttons[i])

		// ボタンクリック時の処理
		tb.buttons[i].OnClick = func() {
			tb.marker.Y = i*32 + 8

			if tb.OnSelect != nil {
				tb.OnSelect(tb.buttons[i].Size)
			}
		}
	}

	return tb
}
