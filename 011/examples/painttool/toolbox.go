package main

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ToolButtonはツール選択ボタンのコントロール
type ToolButton struct {
	parts.ControlBase      // 基本的なコントロール機能
	parts.MouseInteraction // マウス操作（クリックなど）の処理
	parts.Drawable         // 描画機能

	Size int
}

// ToolButton生成
func NewToolButton(x, y, w, h, size int) *ToolButton {
	b := &ToolButton{}
	b.InitControlBase(b, x, y, w, h)
	b.InitMouseInteraction(b)
	b.InitDrawable(b)
	b.OnDraw = b.drawButton
	b.Size = size
	return b
}

// drawButtonはボタンの描画処理を行う
func (b *ToolButton) drawButton(screen *ebiten.Image) {
	gx, gy := b.GetGlobalPos()
	vector.StrokeRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), 2, color.RGBA{0xff, 0xff, 0xff, 0xff}, false)
	vector.FillCircle(screen, float32(gx+b.Width/2), float32(gy+b.Height/2), float32(b.Size/2), color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
}

// ToolMarkerはツール選択マーカーのコントロール
type ToolMarker struct {
	parts.ControlBase
	parts.TextDrawable
}

// ToolMarker生成
func NewToolMarker(x, y, w, h int) *ToolMarker {
	m := &ToolMarker{}
	m.InitControlBase(m, x, y, w, h)
	m.InitTextDrawable(m, "▶", 16, parts.AlignCenter, parts.AlignCenter, 0, 0, color.RGBA{0xff, 0xff, 0x00, 0xff}, false)
	return m
}

// ToolBoxはツール選択ボックスのコントロール
type ToolBox struct {
	parts.ControlBase
	parts.Grouping

	marker   *ToolMarker
	buttons  []*ToolButton
	OnSelect func(int)
}

// ToolBox生成
func NewToolBox(x, y, w, h int) *ToolBox {
	tb := &ToolBox{}
	tb.InitControlBase(tb, x, y, w, h)
	tb.InitGrouping(tb)
	// ToolMarker生成と配置
	tb.marker = NewToolMarker(0, 2*32+8, 16, 16)
	tb.Grouping.AddChild(tb.marker)

	// ToolButton生成と配置
	for i := 0; i < 8; i++ {
		b := NewToolButton(18, i*32, 32, 32, i*2+2)
		tb.AddChild(b)
		tb.buttons = append(tb.buttons, b)

		// ボタンクリック時の処理
		b.OnClick = func() {
			tb.marker.Y = b.Y + 8

			if tb.OnSelect != nil {
				tb.OnSelect(b.Size)
			}
		}
	}

	return tb
}
