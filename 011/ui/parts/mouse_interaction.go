package parts

import "MyProject/ui/input"

// MouseInteractionはマウスやタッチ操作（クリック、ドラッグ、ホバーなど）を処理する機能
// コントロールに埋め込むことで使用する
type MouseInteraction struct {
	Control Control
	touch   input.Touch

	OnPress        func()
	OnClick        func()
	OnHoverStart   func()
	OnHoverEnd     func()
	OnRightPress   func()
	OnRightRelease func()
	OnDragStart    func(x, y int)
	OnDrag         func(x, y int)
	OnDragEnd      func()
	OnRelease      func()
	OnRepeat       func() // オートリピート（押しっぱなしで連続実行）

	RepeatDelay    int // リピート開始までの待機フレーム数
	RepeatInterval int // リピート間隔のフレーム数
	pressDuration  int // 押されている時間

	IsPressed   bool // 押されている状態
	IsDragging  bool // ドラッグ中の状態
	IsHovering  bool // カーソルが上に乗っている状態(押されていない)
	IsMouseOver bool // カーソルが上に乗っている状態(押されていても)

	wasHovering bool // 前フレームのホバー状態
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

// handleInputFunctionは入力イベントを処理する
// クリック開始、ドラッグ開始、ホバー状態の判定を行う
func (m *MouseInteraction) handleInputFunction(t input.Touch) bool {
	if !m.Control.GetControlBase().Visible {
		m.touch = nil
		m.IsHovering = false
		m.IsMouseOver = false
		m.IsPressed = false
		m.IsDragging = false

		if m.wasHovering {
			if m.OnHoverEnd != nil {
				m.OnHoverEnd()
			}
			m.wasHovering = false
		}
		return false
	}

	handle := false
	cb := m.Control.GetControlBase()
	gx, gy := cb.GetGlobalPos()

	m.IsPressed = false
	m.IsDragging = false
	m.IsHovering = false
	m.IsMouseOver = false

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

			m.IsMouseOver = true

			x, y := input.GetMouseTouch().Pos()
			// 座標判定
			if gx <= x && x < gx+cb.Width && gy <= y && y < gy+cb.Height {
				// 右クリック検出
				if input.IsRightJustPressed() {
					if m.OnRightPress != nil {
						m.OnRightPress()
					}
				}
				if input.IsRightJustReleased() {
					if m.OnRightRelease != nil {
						m.OnRightRelease()
					}
				}
			}

			handle = true
		}
	} else if m.touch != nil {
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
				m.IsMouseOver = true
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
				m.IsMouseOver = true
			}

			m.IsDragging = true
		}
	}

	// ホバー状態の変化を検知
	if !m.wasHovering && m.IsHovering {
		if m.OnHoverStart != nil {
			m.OnHoverStart()
		}
	} else if m.wasHovering && !m.IsHovering {
		if m.OnHoverEnd != nil {
			m.OnHoverEnd()
		}
	}
	m.wasHovering = m.IsHovering

	return handle
}

// 現在のマウス/タッチ位置を取得する
func (m *MouseInteraction) GetPosition() (int, int, bool) {
	if m.touch != nil {
		x, y := m.touch.Pos()
		return x, y, true
	}
	return 0, 0, false
}
