package parts

import "MyProject/ui/input"

// MouseInteractionはマウスやタッチ操作（クリック、ドラッグ、プレス）を処理する機能を提供します。
// コントロールに埋め込んで使用します。
type MouseInteraction struct {
	Control Control
	touch   input.Touch

	OnPress      func()
	OnClick      func()
	OnRightClick func() // 右クリック
	OnDragStart  func(x, y int)
	OnDrag       func(x, y int)
	OnDragEnd    func()
	OnRelease    func()
	OnRepeat     func() // オートリピート（押しっぱなしで連続実行）

	RepeatDelay    int // リピート開始までの待機フレーム数
	RepeatInterval int // リピート間隔のフレーム数
	pressDuration  int // 押されている時間

	IsPressed  bool // 押されている状態
	IsDragging bool // ドラッグ中の状態
	IsHovering bool // カーソルが上に乗っている状態(押されていない)
	IsOver     bool // カーソルが上に乗っている状態(押されていても)
}

func NewMouseInteraction(c Control) *MouseInteraction {
	m := &MouseInteraction{}
	m.InitMouseInteraction(c)
	return m
}

func (m *MouseInteraction) InitMouseInteraction(c Control) {
	m.Control = c
	m.RepeatDelay = 30
	m.RepeatInterval = 5

	// コントロールのHandleInput時に呼ばれる関数を登録する
	c.GetControlBase().AddHandleInputFunction(m.handleInputFunction)
}

// handleInputFunctionは入力イベントを処理します。
// クリック開始、ドラッグ開始、ホバー状態の判定を行います。
func (m *MouseInteraction) handleInputFunction(t input.Touch) bool {
	cb := m.Control.GetControlBase()
	gx, gy := cb.GetGlobalPos()

	m.IsPressed = false
	m.IsDragging = false
	m.IsHovering = false
	m.IsOver = false

	if m.touch == nil && t != nil {
		x, y := t.Pos()
		// 座標判定
		if gx <= x && x < gx+cb.Width && gy <= y && y < gy+cb.Height {
			if t.IsJustPressed() { // 今回押された
				m.touch = t

				if m.OnPress != nil {
					m.OnPress()
				}

				if m.OnDragStart != nil {
					m.OnDragStart(x, y)
				}

				m.IsPressed = true
				m.IsDragging = true
			}

			if !t.IsPressed() { // 押されていない
				m.IsHovering = true
			}

			if t.IsJustReleased() { // 今回離された
				if m.OnRelease != nil {
					m.OnRelease()
				}
			}

			m.IsOver = true

			return true
		}
	}

	if m.touch != nil {
		x, y := m.touch.Pos()

		if m.touch.IsJustReleased() { // 今回離された
			m.touch = nil
			m.pressDuration = 0

			if m.OnRelease != nil {
				m.OnRelease()
			}

			if m.OnDragEnd != nil {
				m.OnDragEnd()
			}

			if gx <= x && x < gx+cb.Width && gy <= y && y < gy+cb.Height {
				if m.OnClick != nil {
					m.OnClick()
				}

				m.IsHovering = true
				m.IsOver = true
			}
		} else { // 継続して押されている
			m.pressDuration++

			if m.OnRepeat != nil {
				if m.pressDuration == 1 {
					m.OnRepeat()
				} else if m.pressDuration > m.RepeatDelay {
					if (m.pressDuration-m.RepeatDelay)%m.RepeatInterval == 0 {
						m.OnRepeat()
					}
				}
			}

			if m.OnDrag != nil {
				// ドラッグ動作：範囲外に出ても継続し、移動先を通知する
				m.OnDrag(x, y)
			}

			if gx <= x && x < gx+cb.Width && gy <= y && y < gy+cb.Height {
				m.IsPressed = true
				m.IsOver = true
			}

			m.IsDragging = true
		}
	}

	x, y := input.GetMouseTouch().Pos()
	// 座標判定
	if gx <= x && x < gx+cb.Width && gy <= y && y < gy+cb.Height {
		// 右クリック検出
		if input.IsRightJustReleased() {
			if m.OnRightClick != nil {
				m.OnRightClick()
			}
		}
	}

	return false
}

// 現在のマウス/タッチ位置を取得する
func (m *MouseInteraction) GetPosition() (int, int, bool) {
	if m.touch != nil {
		x, y := m.touch.Pos()
		return x, y, true
	}
	return 0, 0, false
}
