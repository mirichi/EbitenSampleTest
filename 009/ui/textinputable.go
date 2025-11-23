package ui

import (
	"image"
	"image/color"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TextInputableはテキスト入力用の更新機能
type TextInputable struct {
	Control   Control
	Field     *textinput.Field
	Counter   *int
	isFocused func() bool

	FontSize    int
	AlignX      TextAlign
	Color       color.Color
	CursorColor color.Color

	selectionAnchor bool // Start側が開始点ならtrue
}

// TextInputable生成
func NewTextInputable(c Control, field *textinput.Field, counter *int, fontSize int, alignX TextAlign, color, cursorColor color.Color, isFocused func() bool) *TextInputable {
	u := &TextInputable{
		Control:         c,
		Field:           field,
		Counter:         counter,
		isFocused:       isFocused,
		FontSize:        fontSize,
		AlignX:          alignX,
		Color:           color,
		CursorColor:     cursorColor,
		selectionAnchor: true,
	}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(u.updateFunction)

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetControlBase().AddDrawFunction(u.drawFunction)

	return u
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (u *TextInputable) updateFunction(t TouchInfo) UpdateResult {
	if u.Field.IsFocused() {
		if u.Counter != nil {
			*u.Counter++
		}

		cb := u.Control.GetControlBase()
		ox, oy := cb.GetOwnerPos()
		x, y := ox+cb.X, oy+cb.Y

		handled, err := u.Field.HandleInputWithBounds(image.Rect(x, y, x+cb.Width, y+cb.Height))
		if err != nil {
			return NotConsumed
		}

		if handled {
			return Consumed
		}

		// BackSpace対応
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			text := u.Field.Text()
			selectionStart, selectionEnd := u.Field.Selection()

			if selectionStart != selectionEnd {
				// 選択範囲がある場合は削除
				newText := text[:selectionStart] + text[selectionEnd:]
				u.Field.SetTextAndSelection(newText, selectionStart, selectionStart)
			} else if selectionStart > 0 {
				// 選択範囲がない場合はカーソル直前の文字を削除
				_, l := utf8.DecodeLastRuneInString(text[:selectionStart])
				newText := text[:selectionStart-l] + text[selectionEnd:]
				selectionStart -= l
				u.Field.SetTextAndSelection(newText, selectionStart, selectionStart)
			}
		}

		// Shiftキーが押されているか
		isShiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift)

		// 矢印キー対応
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			text := u.Field.Text()
			selectionStart, selectionEnd := u.Field.Selection()

			if isShiftPressed {
				// 未選択状態もしくはEndがアンカーの場合
				if selectionStart == selectionEnd || !u.selectionAnchor {
					u.selectionAnchor = false

					// Startを左に移動
					if selectionStart > 0 {
						_, l := utf8.DecodeLastRuneInString(text[:selectionStart])
						selectionStart -= l
					}
				} else {
					u.selectionAnchor = true

					// Endを左に移動
					if selectionEnd > 0 {
						_, l := utf8.DecodeLastRuneInString(text[:selectionEnd])
						selectionEnd -= l
					}
				}
			} else {
				// 選択状態の場合、解除
				if selectionStart != selectionEnd {
					selectionEnd = selectionStart
				} else {
					// 左に移動
					if selectionStart > 0 {
						_, l := utf8.DecodeLastRuneInString(text[:selectionStart])
						selectionStart -= l
					}
					selectionEnd = selectionStart
				}
			}

			u.Field.SetSelection(selectionStart, selectionEnd)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			text := u.Field.Text()
			selectionStart, selectionEnd := u.Field.Selection()

			if isShiftPressed {
				// 未選択状態もしくはStartがアンカーの場合
				if selectionStart == selectionEnd || u.selectionAnchor {
					u.selectionAnchor = true

					// Endを右に移動
					if selectionEnd < len(text) {
						selectionEnd = selectionEnd + 1
						for selectionEnd < len(text) && (text[selectionEnd]&0xC0) == 0x80 {
							selectionEnd++
						}
					}
				} else {
					u.selectionAnchor = false

					// Startを右に移動
					selectionStart = selectionStart + 1
					for selectionStart < len(text) && (text[selectionStart]&0xC0) == 0x80 {
						selectionStart++
					}
				}
			} else {
				// 選択状態の場合、解除
				if selectionStart != selectionEnd {
					selectionStart = selectionEnd
				} else {
					// 右に移動
					if selectionStart < len(text) {
						selectionStart = selectionStart + 1
						for selectionStart < len(text) && (text[selectionStart]&0xC0) == 0x80 {
							selectionStart++
						}
					}
					selectionEnd = selectionStart
				}
			}

			u.Field.SetSelection(selectionStart, selectionEnd)
		}
	}

	return NotConsumed
}

// コントロールのDraw時に呼ばれるDrawFunction
func (d *TextInputable) drawFunction(screen *ebiten.Image) {
	f := &text.GoTextFace{
		Source: MplusFaceSource,
		Size:   float64(d.FontSize),
	}

	// 描画テキスト取得
	txt := d.Field.TextForRendering()

	// 描画幅取得
	mw, _ := text.Measure(txt, f, 0)

	// 描画座標算出
	cb := d.Control.GetControlBase()
	ox, oy := cb.GetOwnerPos()
	x, y := float64(cb.X+ox), float64(cb.Y+oy)

	// 横方向整列
	switch d.AlignX {
	case AlignLeft:
		// そのまま
	case AlignCenter:
		x += (float64(cb.Width) - mw) / 2
	case AlignRight:
		x += (float64(cb.Width) - mw)
	}

	// 縦方向整列（中央寄せ固定）
	m := f.Metrics()
	y += (float64(cb.Height - int(m.HAscent) - int(m.HDescent))) / 2

	// 選択範囲の背景描画
	selectionStart, selectionEnd := d.Field.Selection()

	// 日本語入力が開始されたら選択範囲は描画しない
	if selectionStart != selectionEnd && d.Field.UncommittedTextLengthInBytes() == 0 {
		// 選択範囲より前のテキスト幅
		textBeforeSelection := txt[:selectionStart]
		offset, _ := text.Measure(textBeforeSelection, f, 0)

		// 選択範囲のテキスト幅
		selectedText := txt[selectionStart:selectionEnd]
		width, _ := text.Measure(selectedText, f, 0)

		// 選択範囲の背景を描画（青色半透明）
		vector.FillRect(screen, float32(x+offset), float32(y), float32(width), float32(m.HAscent+m.HDescent), color.RGBA{0x00, 0x00, 0xff, 0x40}, false)
	}

	// テキスト描画
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(d.Color)
	text.Draw(screen, txt, f, op)

	// カーソル描画
	if d.isFocused() && *d.Counter%60 < 30 {
		// 選択位置までのテキスト幅を計測
		selectionStart, selectionEnd := d.Field.Selection()

		// CompositionSelectionは未実装？らしく0しか返ってこないのでこっちを使用
		compLen := d.Field.UncommittedTextLengthInBytes()
		if compLen != 0 {
			selectionEnd = selectionStart + compLen
		} else {
			if !d.selectionAnchor {
				selectionEnd = selectionStart
			}
		}

		textBeforeCursor := txt[:selectionEnd]
		cursorMw, _ := text.Measure(textBeforeCursor, f, 0)

		cursorX := x + cursorMw

		// カーソル線を描画
		vector.StrokeLine(screen, float32(cursorX), float32(y), float32(cursorX), float32(y+m.HAscent+m.HDescent), 1, d.CursorColor, false)
	}
}
