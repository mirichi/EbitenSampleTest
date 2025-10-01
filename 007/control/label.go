package control

import (
	"image/color"

	"myproject/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Label struct {
	ui.ControlBase // これを埋め込むとコントロールとして扱える
	Label          string
	Size           float64
	Adjust         ui.TextAdjust
	color.Color
}

func NewLabel(x, y, w, h int, l string, sz float64, adj ui.TextAdjust, o ui.Control) *Label {
	return &Label{
		ControlBase: ui.ControlBase{X: x, Y: y, W: w, H: h, Owner: o},
		Label:       l,
		Size:        sz,
		Adjust:      adj,
		Color:       color.White,
	}
}

// 描画
func (l *Label) Draw(screen *ebiten.Image) {
	tx, ty := l.GetGlobalPos()
	ox, oy := float64(tx), float64(ty)

	mw, _ := text.Measure(l.Label, &text.GoTextFace{
		Source: ui.MplusFaceSource,
		Size:   l.Size,
	}, 0)

	switch l.Adjust {
	case ui.AdjustCenter:
		ox += (float64(l.W) - mw) / 2

	case ui.AdjustRight:
		ox += (float64(l.W) - mw)

	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(ox), float64(oy))
	op.ColorScale.ScaleWithColor(l.Color)
	text.Draw(screen, l.Label, &text.GoTextFace{
		Source: ui.MplusFaceSource,
		Size:   l.Size,
	}, op)
}

func (l *Label) ProcessTouch(t ui.TouchInfo) bool {
	return false
}

func (l *Label) Update() {
}
