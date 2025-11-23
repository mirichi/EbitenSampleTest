package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	resizeMargin = 10
)

// Resizableはコントロールのサイズ変更操作を可能にする機能
type Resizable struct {
	Control Control
	TouchInfo
	mode     int    // 12346789でサイズ変更中を表す
	OnResize func() // リサイズ時に呼ぶ関数
}

// Resizalbe生成
// fはリサイズ時に実行される関数
func NewResizable(c Control) *Resizable {
	r := &Resizable{
		Control: c,
	}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(r.updateFunction)

	// コントロールのDraw時に呼ばれる関数を登録する
	c.GetControlBase().AddDrawFunction(r.drawFunction)

	return r
}

// 8方向のリサイズ判定
func (r *Resizable) judgeResize(mx, my int) int {
	c := r.Control.GetControlBase()
	ox, oy := c.GetOwnerPos()
	x := ox + c.X
	y := oy + c.Y
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

// コントロールのUpdate時に呼ばれるUpdateFunction
func (r *Resizable) updateFunction(t TouchInfo) UpdateResult {
	var result UpdateResult = NotConsumed

	if r.TouchInfo == nil { // 押されていないとき
		if t != nil {
			if t.IsJustPressed() { // 今回押された
				mx, my := t.Pos()
				r.mode = r.judgeResize(mx, my)
				if r.mode > 0 {
					r.TouchInfo = t
					result = StopPropagation
				}
			} else {
				mx, my := t.Pos()
				r.mode = r.judgeResize(mx, my)
				if r.mode > 0 {
					r.updateCursor(r.mode)
					result = StopPropagation
				}
			}
		}
	} else { // 押されているとき
		if r.TouchInfo.IsJustReleased() { // 今回離された
			r.TouchInfo = nil
		} else { // 継続して押されている
			r.updateCursor(r.mode)

			x, y := r.TouchInfo.Pos()
			bx, by := r.TouchInfo.OldPos()
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

			// サイズ0にならない程度の雑い対策
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

			result = Consumed
		}
	}

	return result
}

// リサイズマーカー描画
func (r *Resizable) drawFunction(screen *ebiten.Image) {
	// マーカー描画は廃止
}
