package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TettBoxは文字入力と描画をするコントロール
type TextBox struct {
	*parts.ControlBase
	*parts.MouseInteraction
	*parts.Drawable
	*parts.TextInputable
	*parts.Focusable

	field   textinput.Field
	counter int
}

// TetxBox生成
func NewTextBox(x, y, w, h int, text string, fontSize int) *TextBox {
	tb := &TextBox{}

	tb.ControlBase = parts.NewControlBase(tb, x, y, w, h)

	// 初期テキスト
	tb.field.SetTextAndSelection(text, len(text), len(text))

	tb.MouseInteraction = parts.NewMouseInteraction(tb)

	// マウス押下時にフォーカス取得と範囲選択開始
	tb.OnPress = func() {
		tb.Focus()
		if x, y, ok := tb.MouseInteraction.GetPosition(); ok {
			tb.TextInputable.StartSelection(x, y)
		}
	}

	// ドラッグ中は選択範囲を更新
	tb.OnDrag = func(x, y int) {
		tb.TextInputable.UpdateSelection(x, y)
	}

	tb.Drawable = parts.NewDrawable(tb, tb.drawTextBox)
	tb.TextInputable = parts.NewTextInputable(tb, &tb.field, &tb.counter, fontSize, parts.AlignLeft, color.White, color.White, func() bool { return tb.Focused })
	tb.Focusable = parts.NewFocusable(tb)
	tb.OnFocus = func() {
		tb.field.Focus()
	}
	tb.OnBlur = func() {
		tb.field.Blur()
	}

	return tb
}

func (tb *TextBox) drawTextBox(screen *ebiten.Image) {
	gx, gy := tb.GetGlobalPos()
	x := float32(gx)
	y := float32(gy)

	// 背景と枠線
	w := float32(tb.Width)
	h := float32(tb.Height)
	vector.FillRect(screen, x, y, w, h, color.RGBA{0x60, 0x60, 0x60, 0xff}, false)

	// 通常時の枠線（フォーカス枠はFocusableが描画）
	if !tb.Focused {
		vector.StrokeRect(screen, x, y, w, h, 2, color.RGBA{0x80, 0x80, 0x80, 0xff}, false)
	}
}
