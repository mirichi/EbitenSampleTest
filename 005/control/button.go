package control

import (
	"image/color"

	"myproject/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	ui.ControlBase // これを埋め込むとコントロールとして扱える
	Label          string
	Size           float64
	Adjust         ui.TextAdjust
	color.Color
}

func NewButton(x, y, w, h int, l string, sz float64, adj ui.TextAdjust, o ui.Control, f func()) *Button {
	return &Button{
		ControlBase: ui.ControlBase{X: x, Y: y, W: w, H: h, Owner: o, TapFunc: f},
		Label:       l,
		Size:        sz,
		Adjust:      adj,
		Color:       color.White,
	}
}

// 描画
func (b *Button) Draw(screen *ebiten.Image) {
	tx, ty := b.GetGlobalPos()
	ox, oy := float64(tx), float64(ty)
	vector.StrokeRect(screen, float32(ox), float32(oy), float32(b.W), float32(b.H), 3, color.RGBA{0xff, 0xff, 0xff, 0xff}, false)

	mw, _ := text.Measure(b.Label, &text.GoTextFace{
		Source: ui.MplusFaceSource,
		Size:   b.Size,
	}, 0)

	switch b.Adjust {
	case ui.AdjustCenter:
		ox += (float64(b.W-10) - mw) / 2

	case ui.AdjustRight:
		ox += (float64(b.W-10) - mw)

	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(ox+5), float64(oy+5))
	op.ColorScale.ScaleWithColor(b.Color)
	text.Draw(screen, b.Label, &text.GoTextFace{
		Source: ui.MplusFaceSource,
		Size:   b.Size,
	}, op)
}

// 状態更新
func (b *Button) Update() {
	b.ControlBase.Update()

	// タッチ中
	if b.TouchInfo != nil {
		b.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
	} else {
		b.Color = color.White
	}
}
