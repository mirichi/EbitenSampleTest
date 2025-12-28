package parts

import (
	"MyProject/ui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	resizeMargin = 7
)

// ResizableはマウスドラッグによるControlのサイズ変更機能を提供する
// Controlの端をドラッグすることでサイズ変更が可能になる
type Resizable struct {
	Control  Control
	touch    input.Touch
	mode     int    // 12346789でサイズ変更中を表す
	OnResize func() // リサイズ時に呼ぶ関数

	// ドラッグ開始時の状態
	startMX, startMY float64
	startX, startY   float64
	startW, startH   float64
}

func NewResizable(c Control) *Resizable {
	r := &Resizable{}
	r.InitResizable(c)
	return r
}

func (r *Resizable) InitResizable(c Control) {
	r.Control = c

	// ControlのHandleInput時に呼ばれる関数を登録する
	c.GetControlBase().AddHandleInputFunction(r.handleInputFunction)
}

// 8方向のリサイズ判定
func (r *Resizable) judgeResize(mx, my float64) int {
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

// ControlのHandleInput時に呼ばれるHandleInputFunction
func (r *Resizable) handleInputFunction(t input.Touch) bool {
	if r.touch == nil && t != nil {
		if t.IsJustPressed() { // 今回押された
			mx, my := t.Pos()
			r.mode = r.judgeResize(mx, my)
			if r.mode > 0 {
				r.touch = t
				// ドラッグ開始時の状態を保存
				r.startMX, r.startMY = mx, my
				c := r.Control.GetControlBase()
				r.startX, r.startY = c.X, c.Y
				r.startW, r.startH = c.Width, c.Height
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

	if r.touch != nil {
		if r.touch.IsJustReleased() { // 今回離された
			r.mode = 0
			r.touch = nil
		} else { // 継続して押されている
			// カーソルは継続して変更しないと自動的に元に戻される
			r.updateCursor(r.mode)

			mx, my := r.touch.Pos()
			dx := mx - r.startMX
			dy := my - r.startMY
			c := r.Control.GetControlBase()

			// 計算用の一時変数
			newX, newY := r.startX, r.startY
			newW, newH := r.startW, r.startH

			// 左上、左、左下
			if r.mode == 7 || r.mode == 4 || r.mode == 1 {
				newW = r.startW - dx
				newX = r.startX + dx
			}

			// 左上、上、右上
			if r.mode == 7 || r.mode == 8 || r.mode == 9 {
				newH = r.startH - dy
				newY = r.startY + dy
			}

			// 右上、右、右下
			if r.mode == 9 || r.mode == 6 || r.mode == 3 {
				newW = r.startW + dx
			}

			// 左下、下、右下
			if r.mode == 1 || r.mode == 2 || r.mode == 3 {
				newH = r.startH + dy
			}

			// 最小サイズの制限
			if newW < 10 {
				// 左側を動かしている場合、X座標を調整して右端を固定する
				if r.mode == 7 || r.mode == 4 || r.mode == 1 {
					newX = r.startX + r.startW - 10
				}
				newW = 10
			}
			if newH < 10 {
				// 上側を動かしている場合、Y座標を調整して下端を固定する
				if r.mode == 7 || r.mode == 8 || r.mode == 9 {
					newY = r.startY + r.startH - 10
				}
				newH = 10
			}

			// 適用
			c.X = newX
			c.Y = newY
			c.Width = newW
			c.Height = newH

			if r.OnResize != nil {
				r.OnResize()
			}
		}
	}

	return false
}
