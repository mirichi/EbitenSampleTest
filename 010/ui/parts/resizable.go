package parts

import (
	"MyProject/ui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	resizeMargin = 7
)

// Resizableはマウスドラッグによるコントロールのサイズ変更機能を提供します。
// コントロールの端をドラッグすることでリサイズが可能になります。
type Resizable struct {
	Control  Control
	touch    input.Touch
	mode     int    // 12346789でサイズ変更中を表す
	OnResize func() // リサイズ時に呼ぶ関数
}

// Resizalbe生成
// fはリサイズ時に実行される関数
func NewResizable(c Control) *Resizable {
	r := &Resizable{
		Control: c,
	}

	// コントロールのHandleInput時に呼ばれる関数を登録する
	c.GetControlBase().AddHandleInputFunction(r.handleInputFunction)

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(r.updateFunction)

	return r
}

// 8方向のリサイズ判定
func (r *Resizable) judgeResize(mx, my int) int {
	c := r.Control.GetControlBase()
	gx, gy := c.GetGlobalPos()
	x := gx
	y := gy
	w := c.Width
	h := c.Height

	// 判定用フラグ
	inLeft := x-resizeMargin <= mx && mx <= x+resizeMargin
	inRight := x+w-resizeMargin <= mx && mx <= x+w+resizeMargin
	inTop := y-resizeMargin <= my && my <= y+resizeMargin
	inBottom := y+h-resizeMargin <= my && my <= y+h+resizeMargin

	inX := x-resizeMargin <= mx && mx <= x+w+resizeMargin
	inY := y-resizeMargin <= my && my <= y+h+resizeMargin

	mode := 0

	// 範囲外なら0
	if !inX || !inY {
		return 0
	}

	if inLeft && inTop {
		mode = 7
	} else if inRight && inTop {
		mode = 9
	} else if inLeft && inBottom {
		mode = 1
	} else if inRight && inBottom {
		mode = 3
	} else if inTop {
		mode = 8
	} else if inBottom {
		mode = 2
	} else if inLeft {
		mode = 4
	} else if inRight {
		mode = 6
	}

	return mode
}

func (r *Resizable) updateCursor(m int) {
	switch m {
	case 1, 9: // 左下、右上
		RequestCursorShape(ebiten.CursorShapeNESWResize)
	case 7, 3: // 左上、右下
		RequestCursorShape(ebiten.CursorShapeNWSEResize)
	case 4, 6: // 左、右
		RequestCursorShape(ebiten.CursorShapeEWResize)
	case 8, 2: //上、下
		RequestCursorShape(ebiten.CursorShapeNSResize)
	}
}

// コントロールのHandleInput時に呼ばれるHandleInputFunction
func (r *Resizable) handleInputFunction(t input.Touch) bool {
	if r.touch == nil {
		if t.IsJustPressed() { // 今回押された
			mx, my := t.Pos()
			r.mode = r.judgeResize(mx, my)
			if r.mode > 0 {
				r.touch = t
				return true
			}
		}

		if !t.IsPressed() { // 押されていない
			mx, my := t.Pos()
			mode := r.judgeResize(mx, my)
			if mode > 0 {
				r.updateCursor(mode)
				return true
			}
		}
	}

	return false
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (r *Resizable) updateFunction() {
	if r.touch != nil {
		if r.touch.IsJustReleased() { // 今回離された
			r.mode = 0
			r.touch = nil
		} else { // 継続して押されている
			// カーソルは継続して変更しないと自動的に元に戻される
			r.updateCursor(r.mode)

			x, y := r.touch.Pos()
			bx, by := r.touch.OldPos()
			c := r.Control.GetControlBase()

			// 左上、左、左下
			if r.mode == 7 || r.mode == 4 || r.mode == 1 {
				c.X += x - bx
				c.Width -= x - bx
			}

			// 左上、上、右上
			if r.mode == 7 || r.mode == 8 || r.mode == 9 {
				c.Y += y - by
				c.Height -= y - by
			}

			// 右上、右、右下
			if r.mode == 9 || r.mode == 6 || r.mode == 3 {
				c.Width += x - bx
			}

			// 左下、下、右下
			if r.mode == 1 || r.mode == 2 || r.mode == 3 {
				c.Height += y - by
			}

			// 最小サイズの制限
			// コントロールが消失したり操作不能になるのを防ぐため、最小サイズを確保します。
			if r.mode > 0 {
				if c.Width < 10 {
					c.Width = 10
				}
				if c.Height < 10 {
					c.Height = 10
				}
				if r.OnResize != nil {
					r.OnResize()
				}
			}
		}
	}
}
