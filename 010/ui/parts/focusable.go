package parts

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Focusable struct {
	Control Control
	Focused bool

	OnFocus func()
	OnBlur  func()

	// フォーカス枠の設定
	DrawFocusBorder  bool        // フォーカス枠を描画するかどうか
	FocusBorderColor color.Color // フォーカス枠の色
	FocusBorderWidth float32     // フォーカス枠の太さ
}

type FocusableInterface interface {
	Focus()
	Blur()
}

var FocusedControl FocusableInterface
var FlameFocus bool

func NewFocusable(c Control) *Focusable {
	f := &Focusable{
		Control:          c,
		DrawFocusBorder:  true,
		FocusBorderColor: color.RGBA{0xD0, 0xD0, 0xD0, 0xff},
		FocusBorderWidth: 2,
	}

	// フォーカス枠を描画する関数を登録
	c.GetControlBase().AddDrawFunction(f.drawFocusBorder)

	return f
}

func (f *Focusable) Focus() {
	if FocusedControl != f.Control.(FocusableInterface) {
		if FocusedControl != nil {
			FocusedControl.Blur()
		}
		f.Focused = true
		FocusedControl = f.Control.(FocusableInterface)
		if f.OnFocus != nil {
			f.OnFocus()
		}
	}
	FlameFocus = true
}

func (f *Focusable) Blur() {
	f.Focused = false
	FocusedControl = nil
	if f.OnBlur != nil {
		f.OnBlur()
	}
}

// drawFocusBorder はフォーカスがある場合に枠を描画します。
func (f *Focusable) drawFocusBorder(screen *ebiten.Image) {
	if !f.DrawFocusBorder || !f.Focused {
		return
	}

	cb := f.Control.GetControlBase()
	gx, gy := cb.GetGlobalPos()
	vector.StrokeRect(screen,
		float32(gx), float32(gy),
		float32(cb.Width), float32(cb.Height),
		f.FocusBorderWidth, f.FocusBorderColor, false)
}
