package ui

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
	control *ControlBase
	Text    string
	Size    int
	AlignX  TextAlign
	AlignY  TextAlign
	OffsetX int
	OffsetY int
	Color   color.Color
	Shadow  bool
}

// TextDrawable生成
func NewTextDrawable(c *ControlBase, text string, size int, alignX, alignY TextAlign, offsetX, offsetY int, color color.Color, shadow bool) *TextDrawable {
	td := &TextDrawable{
		control: c,
		Text:    text,
		Size:    size,
		AlignX:  alignX,
		AlignY:  alignY,
		OffsetX: offsetX,
		OffsetY: offsetY,
		Color:   color,
		Shadow:  shadow,
	}

	// コントロールのDraw時に呼ばれる関数を登録する
	c.AddDrawFunction(td.DrawFunction)

	return td
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *TextDrawable) DrawFunction(screen *ebiten.Image) {
	f := &text.GoTextFace{
		Source: MplusFaceSource,
		Size:   float64(d.Size),
	}

	// 描画幅取得
	mw, _ := text.Measure(d.Text, f, 0)

	// 描画座標算出
	ox, oy := d.control.GetOwnerPos()
	x, y := float64(d.control.X+ox), float64(d.control.Y+oy)

	// 横方向整列
	switch d.AlignX {
	case AlignLeft:
		x += float64(d.OffsetX)
	case AlignCenter:
		x += (float64(d.control.Width)-mw)/2 + float64(d.OffsetX)
	case AlignRight:
		x += (float64(d.control.Width+d.OffsetX) - mw)
	}

	// 縦方向整列
	m := f.Metrics()
	switch d.AlignY {
	case AlignUpper:
		y += float64(d.OffsetY)
	case AlignCenter:
		y += (float64(d.control.Height-int(m.HAscent)-int(m.HDescent)))/2 + float64(d.OffsetY)
	case AlignBottom:
		y += (float64(d.control.Height + d.OffsetY - int(m.HAscent) - int(m.HDescent)))
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
