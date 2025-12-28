package main

import (
	"MyProject/ui"
	"image/color"
)

// ColorPanelは色選択パネル
type ColorPanel struct {
	ui.GroupingControl

	marker   ui.Label
	Buttons  [8]ui.InteractiveControl
	OnSelect func(color.Color)
}

func NewColorPanel(x, y, w, h float64) *ColorPanel {
	cp := &ColorPanel{}
	cp.InitGroupingControl(x, y, w, h)

	// ColorMarker初期化
	cp.marker.InitLabel(32, 7*32+8, 16, 16, "◀", 16)
	cp.AddChild(&cp.marker)

	// ColorButton初期化
	cols := []color.Color{color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0xff, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff}}
	for i, c := range cols {
		cp.Buttons[i].InitInteractiveControl(0, float64(i)*32, 32, 32)
		cp.AddChild(&cp.Buttons[i])
		cp.Buttons[i].BackColor = c
		cp.Buttons[i].OnPress = func() {
			cp.marker.Y = float64(i)*32 + 8
			cp.marker.TextColor = c
			if cp.OnSelect != nil {
				cp.OnSelect(c)
			}
		}
	}

	return cp
}
