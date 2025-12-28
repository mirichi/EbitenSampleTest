package ui

import (
	"MyProject/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TextBoxは文字入力と描画をするControl
type TextBox struct {
	InteractiveControl
	parts.TextInputable
	parts.Focusable

	field   textinput.Field
	counter int
}

// TextBox生成
func NewTextBox(x, y, w, h float64, text string, fontSize float64) *TextBox {
	tb := &TextBox{}
	tb.InitTextBox(x, y, w, h, text, fontSize)
	return tb
}

func (tb *TextBox) InitTextBox(x, y, w, h float64, text string, fontSize float64) {
	// 初期テキスト
	tb.field.SetTextAndSelection(text, len(text), len(text))

	tb.InitInteractiveControl(x, y, w, h)
	tb.InitTextInputable(tb, &tb.field, &tb.counter, fontSize, parts.AlignLeft, color.White, color.White, func() bool { return tb.Focused })
	tb.InitFocusable(tb)
	tb.OnDraw = tb.drawTextBox

	tb.OnFocus = func() {
		tb.field.Focus()
	}
	tb.OnBlur = func() {
		tb.field.Blur()
	}

	// マウス押下時にフォーカス取得と範囲選択開始
	tb.OnPress = func() {
		tb.Focus()
		if x, y, ok := tb.MouseInteraction.GetPosition(); ok {
			tb.TextInputable.StartSelection(x, y)
		}
	}

	// ドラッグ中は選択範囲を更新
	tb.OnDrag = func(x, y float64) {
		tb.TextInputable.UpdateSelection(x, y)
	}

}

func (tb *TextBox) drawTextBox(screen *ebiten.Image) {
	gx, gy := tb.GetGlobalPos()
	x := float32(gx)
	y := float32(gy)

	// 背景と枠線
	w := float32(tb.Width)
	h := float32(tb.Height)
	vector.FillRect(screen, x, y, w, h, color.RGBA{0x60, 0x60, 0x60, 0xff}, false)

	// 枠線（フォーカス枠はFocusableが描画）
	vector.StrokeRect(screen, x, y, w-1, h-1, 2, color.RGBA{0x80, 0x80, 0x80, 0xff}, false)
}
