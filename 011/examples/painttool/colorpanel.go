package main

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ColorButtonは色選択ボタンのコントロール
type ColorButton struct {
	parts.ControlBase      // 基本的なコントロール機能
	parts.MouseInteraction // マウス操作（クリックなど）の処理
	parts.Drawable         // 描画機能
	Color                  color.Color
}

// ColorButton生成
func NewColorButton(x, y, w, h int, color color.Color) *ColorButton {
	b := &ColorButton{}
	b.InitControlBase(b, x, y, w, h)
	b.InitMouseInteraction(b)
	b.InitDrawable(b, b.drawButton)
	b.Color = color
	return b
}

// drawButtonはボタンの描画処理を行う
func (b *ColorButton) drawButton(screen *ebiten.Image) {
	gx, gy := b.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), b.Color, false)
}

// ColorMarkerは色選択マーカーのコントロール
type ColorMarker struct {
	parts.ControlBase
	parts.TextDrawable
}

// ColorMarker生成
func NewColorMarker(x, y, w, h int) *ColorMarker {
	m := &ColorMarker{}
	m.InitControlBase(m, x, y, w, h)
	m.InitTextDrawable(m, "◀", 16, parts.AlignCenter, parts.AlignCenter, 0, 0, color.RGBA{0xff, 0xff, 0x00, 0xff}, false)
	return m
}

// ColorPanelは色選択パネルのコントロール
type ColorPanel struct {
	parts.ControlBase
	parts.Grouping

	marker   *ColorMarker
	buttons  []*ColorButton
	OnSelect func(color.Color)
}

func NewColorPanel(x, y, w, h int) *ColorPanel {
	cp := &ColorPanel{}
	cp.InitControlBase(cp, x, y, w, h)
	cp.InitGrouping(&cp.ControlBase)

	// ColorMarker生成と配置
	cp.marker = NewColorMarker(32, 7*32+8, 16, 16)
	cp.Grouping.AddChild(cp.marker)

	// ColorButton生成と配置
	c := []color.Color{color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0xff, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff}}
	for i := 0; i < 8; i++ {
		b := NewColorButton(0, i*32, 32, 32, c[i])
		cp.AddChild(b)
		cp.buttons = append(cp.buttons, b)

		// ボタンクリック時の処理
		b.OnClick = func() {
			cp.marker.Y = b.Y + 8
			cp.marker.Color = b.Color
			if cp.OnSelect != nil {
				cp.OnSelect(b.Color)
			}
		}
	}

	return cp
}
