package main

import (
	"MyProject/ui"
	"MyProject/ui/parts"
	"image/color"
)

// ColorMarkerは色選択マーカー
type ColorMarker struct {
	parts.WidgetBase
	parts.TextDrawable
}

// ColorMarker初期化
func (m *ColorMarker) InitColorMarker(g parts.AddChilder, x, y, w, h int) {
	m.InitWidgetBase(m, x, y, w, h)
	m.InitTextDrawable(m, "◀", 16, parts.AlignCenter, parts.AlignCenter, 0, 0, color.RGBA{0xff, 0xff, 0x00, 0xff}, false)
	if g != nil {
		g.AddChild(m)
	}
}

// ColorPanelは色選択パネル
type ColorPanel struct {
	ui.BlankWidget

	marker   ColorMarker
	buttons  [8]ui.InteractiveWidget
	OnSelect func(color.Color)
}

func NewColorPanel(x, y, w, h int) *ColorPanel {
	cp := &ColorPanel{}
	cp.InitBlankWidget(nil, x, y, w, h)

	// ColorMarker初期化
	cp.marker.InitColorMarker(&cp.BlankWidget, 32, 7*32+8, 16, 16)

	// ColorButton初期化
	cols := []color.Color{color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0xff, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff}}
	for i, c := range cols {
		cp.buttons[i].InitInteractiveWidget(cp, 0, i*32, 32, 32)
		cp.buttons[i].Color = c
		cp.buttons[i].OnPress = func() {
			cp.marker.Y = i*32 + 8
			cp.marker.TextColor = c
			if cp.OnSelect != nil {
				cp.OnSelect(c)
			}
		}
	}

	return cp
}
