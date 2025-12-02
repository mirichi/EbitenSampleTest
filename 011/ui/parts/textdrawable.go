package parts

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// テキストの整列型
type TextAlign int

const (
	AlignLeft   TextAlign = 0
	AlignCenter TextAlign = 1
	AlignRight  TextAlign = 2
	AlignUpper  TextAlign = 3
	AlignBottom TextAlign = 4
)

// とりあえずのフォント
var MplusFaceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	MplusFaceSource = s
}

// TextDrawableはテキストを描画する機能
type TextDrawable struct {
	Control Control
	Text    string
	Size    int
	AlignX  TextAlign
	AlignY  TextAlign
	OffsetX int
	OffsetY int
	Color   color.Color
	Shadow  bool
}

func (d *TextDrawable) InitTextDrawable(c Control, text string, size int, alignX, alignY TextAlign, offsetX, offsetY int, color color.Color, shadow bool) {
	d.Control = c
	d.Text = text
	d.Size = size
	d.AlignX = alignX
	d.AlignY = alignY
	d.OffsetX = offsetX
	d.OffsetY = offsetY
	d.Color = color
	d.Shadow = shadow

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetControlBase().AddDrawFunction(d.drawFunction)
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *TextDrawable) drawFunction(screen *ebiten.Image) {
	f := &text.GoTextFace{
		Source: MplusFaceSource,
		Size:   float64(d.Size),
	}
	// 描画幅取得
	mw, _ := text.Measure(d.Text, f, 0)

	// 描画座標算出
	cb := d.Control.GetControlBase()
	gx, gy := cb.GetGlobalPos()
	x, y := float64(gx), float64(gy)

	// 横方向整列
	switch d.AlignX {
	case AlignLeft:
		x += float64(d.OffsetX)
	case AlignCenter:
		x += (float64(cb.Width)-mw)/2 + float64(d.OffsetX)
	case AlignRight:
		x += (float64(cb.Width+d.OffsetX) - mw)
	}

	// 縦方向整列
	m := f.Metrics()
	switch d.AlignY {
	case AlignUpper:
		y += float64(d.OffsetY)
	case AlignCenter:
		y += (float64(cb.Height-int(m.HAscent)-int(m.HDescent)))/2 + float64(d.OffsetY)
	case AlignBottom:
		y += (float64(cb.Height + d.OffsetY - int(m.HAscent) - int(m.HDescent)))
	}

	// 影描画
	if d.Shadow {
		op := &text.DrawOptions{}
		op.GeoM.Translate(x+2, y+2)
		op.ColorScale.ScaleWithColor(color.Black)
		text.Draw(screen, d.Text, f, op)
	}
	// テキスト描画
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(d.Color)
	text.Draw(screen, d.Text, f, op)
}
