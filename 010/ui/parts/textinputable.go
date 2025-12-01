package parts

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

// TextInputableはテキスト入力機能を提供する構造体です。
// テキストボックスなどのコントロールに埋め込んで使用します。
type TextInputable struct {
	Control   Control          // 親コントロール
	Field     *textinput.Field // Ebitengineのテキスト入力フィールド
	Counter   *int             // カーソル点滅用カウンタ
	isFocused func() bool      // フォーカス状態判定関数

	FontSize    int         // フォントサイズ
	AlignX      TextAlign   // 横方向の配置 (AlignLeft, AlignCenter, AlignRight)
	Color       color.Color // テキスト色
	CursorColor color.Color // カーソル色

	selectionAnchor   bool // 選択範囲のアンカー位置フラグ (true: Start側が固定, false: End側が固定)
	selectionStartIdx int  // ドラッグ選択開始時の文字インデックス
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

// updateFunctionは毎フレームの更新処理を行います。
// テキスト入力、カーソル移動、選択範囲の操作を処理します。
func (u *TextInputable) updateFunction() {
	if u.Field.IsFocused() {
		if u.Counter != nil {
			*u.Counter++
		}

		cb := u.Control.GetControlBase()
		gx, gy := cb.GetGlobalPos()
		x, y := gx, gy

		handled, err := u.Field.HandleInputWithBounds(image.Rect(x, y, x+cb.Width, y+cb.Height))
		if err != nil {
			return
		}

		if handled {
			return
		}

		// BackSpaceキーの処理
		// 選択範囲がある場合は範囲削除、ない場合はカーソル直前の文字を削除します。
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

		// 左矢印キーの処理
		// Shiftキーとの組み合わせで選択範囲の変更、単独押下でカーソル移動を行います。
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

		// 右矢印キーの処理
		// Shiftキーとの組み合わせで選択範囲の変更、単独押下でカーソル移動を行います。
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
}

// drawFunctionはテキストとカーソル、選択範囲の描画を行います。
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
	gx, gy := cb.GetGlobalPos()
	x, y := float64(gx), float64(gy)

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

// getIndexFromPosは指定された座標(x, y)に対応するテキストの文字インデックスを取得します。
// クリック位置からカーソル位置を特定するために使用されます。
func (u *TextInputable) getIndexFromPos(x, y int) int {
	f := &text.GoTextFace{
		Source: MplusFaceSource,
		Size:   float64(u.FontSize),
	}

	// 描画テキスト取得
	txt := u.Field.TextForRendering()

	// 描画幅取得
	mw, _ := text.Measure(txt, f, 0)

	// 描画座標算出
	// 描画座標算出
	cb := u.Control.GetControlBase()
	gx, _ := cb.GetGlobalPos()
	baseX := float64(gx)

	// 横方向整列
	switch u.AlignX {
	case AlignLeft:
		// そのまま
	case AlignCenter:
		baseX += (float64(cb.Width) - mw) / 2
	case AlignRight:
		baseX += (float64(cb.Width) - mw)
	}

	// クリック位置の相対座標
	relX := float64(x) - baseX

	// テキストの範囲外（左側）なら先頭
	if relX < 0 {
		return 0
	}

	// テキストの範囲外（右側）なら末尾
	if relX > mw {
		return len(txt)
	}

	// 文字ごとに幅を計測して位置を特定
	var currentX float64 = 0
	for i, r := range txt {
		rw, _ := text.Measure(string(r), f, 0)
		// 文字の中心を超えたら次の文字とする
		if relX < currentX+rw/2 {
			return i
		}
		currentX += rw
	}

	return len(txt)
}

// 選択開始（クリック時）
func (u *TextInputable) StartSelection(x, y int) {
	idx := u.getIndexFromPos(x, y)
	u.Field.SetSelection(idx, idx)
	u.selectionStartIdx = idx
	u.selectionAnchor = true
}

// 選択更新（ドラッグ時）
func (u *TextInputable) UpdateSelection(x, y int) {
	idx := u.getIndexFromPos(x, y)

	start := u.selectionStartIdx
	end := idx

	if start > end {
		start, end = end, start
		u.selectionAnchor = false // End側（元々のStart）が固定点、Start側（現在位置）が動く点
	} else {
		u.selectionAnchor = true // Start側（元々のStart）が固定点、End側（現在位置）が動く点
	}

	u.Field.SetSelection(start, end)
}
