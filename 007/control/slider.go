package control

import (
	"fmt"
	"image/color"

	"myproject/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Slider struct {
	ui.ControlBase // これを埋め込むとコントロールとして扱える
	Label          string
	Size           float64
	Adjust         ui.TextAdjust
	color.Color
}

func NewSlider(x, y, w, h int, l string, sz float64, adj ui.TextAdjust, o ui.Control, f func()) *Slider {
	return &Slider{
		ControlBase: ui.ControlBase{X: x, Y: y, W: w, H: h, Owner: o, SlideFunc: f},
		Label:       l,
		Size:        sz,
		Adjust:      adj,
		Color:       color.White,
	}
}

// 描画
func (b *Slider) Draw(screen *ebiten.Image) {
	vector.StrokeLine(screen, float32(0), float32(b.Y+b.H/2), float32(640), float32(b.Y+b.H/2), 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)

	tx, ty := b.GetGlobalPos()
	ox, oy := float64(tx), float64(ty)
	vector.DrawFilledRect(screen, float32(ox), float32(oy), float32(b.W), float32(b.H), color.RGBA{0xd0, 0xd0, 0xd0, 0xff}, true)

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
	op.ColorScale.ScaleWithColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	text.Draw(screen, b.Label, &text.GoTextFace{
		Source: ui.MplusFaceSource,
		Size:   b.Size,
	}, op)
}

// 状態更新
func (b *Slider) Update() {
	b.ControlBase.Update()
}

func (b *Slider) Slide() {
	x, _ := b.TouchInfo.Pos()
	ox, _ := b.TouchInfo.OldPos()
	b.X += x - ox

	if b.X < 0 {
		b.X = 0
	}
	if b.X+b.W > 640 {
		b.X = 640 - b.W
	}

	b.Label = fmt.Sprintf("%d", int(b.X*100/(640-b.W)))
}
