package parts

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Focusable struct {
	Widget  Widget
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

var FocusedWidget FocusableInterface
var FlameFocus bool

func NewFocusable(c Widget) *Focusable {
	f := &Focusable{}
	f.InitFocusable(c)
	return f
}

func (f *Focusable) InitFocusable(c Widget) {
	f.Widget = c
	f.DrawFocusBorder = true
	f.FocusBorderColor = color.RGBA{0xD0, 0xD0, 0xD0, 0xff}
	f.FocusBorderWidth = 2

	// フォーカス枠を描画する関数を登録
	c.GetWidgetBase().AddDrawFunction(f.drawFocusBorder)
}

func (f *Focusable) Focus() {
	if FocusedWidget != f.Widget.(FocusableInterface) {
		if FocusedWidget != nil {
			FocusedWidget.Blur()
		}
		f.Focused = true
		FocusedWidget = f.Widget.(FocusableInterface)
		if f.OnFocus != nil {
			f.OnFocus()
		}
	}
	FlameFocus = true
}

func (f *Focusable) Blur() {
	f.Focused = false
	FocusedWidget = nil
	if f.OnBlur != nil {
		f.OnBlur()
	}
}

// drawFocusBorder はフォーカスがある場合に枠を描画する
func (f *Focusable) drawFocusBorder(screen *ebiten.Image) {
	if !f.DrawFocusBorder || !f.Focused {
		return
	}

	cb := f.Widget.GetWidgetBase()
	gx, gy := cb.GetGlobalPos()
	vector.StrokeRect(screen,
		float32(gx), float32(gy),
		float32(cb.Width), float32(cb.Height),
		f.FocusBorderWidth, f.FocusBorderColor, false)
}
