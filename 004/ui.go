package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var mplusFaceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

type control interface {
	CheckPoint(x, y int) bool
	Press(t TouchInfo)
	Update()
	Draw(screen *ebiten.Image)
}

type Button struct {
	X, Y          float64
	Width, Height float64
	Label         string
	Touch         TouchInfo
	Func          func()
}

func NewButton(x, y, w, h float64, l string, f func()) *Button {
	return &Button{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
		Label:  l,
		Touch:  nil,
		Func:   f,
	}
}

// 描画
func (b *Button) Draw(screen *ebiten.Image) {
	vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), 3, color.RGBA{0xff, 0xff, 0xff, 0xff}, false)

	op := &text.DrawOptions{}
	op.GeoM.Translate(b.X+5, b.Y+5)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, b.Label, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}, op)
}

// 座標(x, y)がButtonの中にあるかどうかをチェックする
func (b *Button) CheckPoint(x, y int) bool {
	return float64(x) >= b.X && float64(x) < b.X+b.Width && float64(y) >= b.Y && float64(y) < b.Y+b.Height
}

// タッチ開始
func (b *Button) Press(t TouchInfo) {
	b.Touch = t
}

// 状態更新
func (b *Button) Update() {
	// タッチ中
	if b.Touch != nil {
		x, y := b.Touch.Pos()
		// 自分にヒットしている
		if b.CheckPoint(x, y) {
			// タッチが離された
			if !b.Touch.IsPressed() {
				b.OnTap()
				b.Touch = nil // タッチ中を解除
			}
		} else {
			// ヒットしていない場合はタッチ中を解除
			b.Touch = nil
		}
	}
}

// タップした
func (b *Button) OnTap() {
	b.Func()
}
