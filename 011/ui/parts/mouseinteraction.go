package parts

import "MyProject/ui/input"

// MouseInteractionはマウスやタッチ操作（クリック、ドラッグ、プレス）を処理する機能を提供します。
// コントロールに埋め込んで使用します。
type MouseInteraction struct {
	Widget Control
	touch  input.Touch

	OnPress      func()
	OnHover      func()
	OnClick      func()
	OnRightClick func() // 右クリック
	OnDragStart  func(x, y int)
	OnDrag       func(x, y int)
	OnDragEnd    func()
	OnPressing   func() // 押されている間毎フレーム呼ばれる
	OnRelease    func()
	OnRepeat     func() // オートリピート（押しっぱなしで連続実行）

	RepeatDelay    int // リピート開始までの待機フレーム数
	RepeatInterval int // リピート間隔のフレーム数
	pressDuration  int // 押されている時間
}

func NewMouseInteraction(c Control) *MouseInteraction {
	m := &MouseInteraction{}
	m.InitMouseInteraction(c)
	return m
}

func (m *MouseInteraction) InitMouseInteraction(c Control) {
	m.Widget = c
	m.RepeatDelay = 30
	m.RepeatInterval = 5

	// コントロールのHandleInput時に呼ばれる関数を登録する
	c.GetControlBase().AddHandleInputFunction(m.handleInputFunction)

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(m.updateFunction)
}

// handleInputFunctionは入力イベントを処理します。
// クリック開始、ドラッグ開始、ホバー状態の判定を行います。
func (m *MouseInteraction) handleInputFunction(t input.Touch) bool {
	if m.touch == nil && t != nil {
		cb := m.Widget.GetControlBase()
		gx, gy := cb.GetGlobalPos()
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
			}

			if !t.IsPressed() { // 押されていない
				if m.OnHover != nil {
					m.OnHover()
				}
			}

			if t.IsJustReleased() { // 今回離された
				if m.OnRelease != nil {
					m.OnRelease()
				}
			}

			return true
		}
	}

	if m.touch != nil {
		if m.touch.IsJustReleased() { // 今回離された
			m.touch = nil
			m.pressDuration = 0

			if m.OnRelease != nil {
				m.OnRelease()
			}

			if m.OnClick != nil {
				m.OnClick()
			}

			if m.OnDragEnd != nil {
				m.OnDragEnd()
			}
		} else { // 継続して押されている
			m.pressDuration++

			if m.OnPressing != nil {
				m.OnPressing()
			}

			if m.OnRepeat != nil {
				if m.pressDuration == 1 {
					m.OnRepeat()
				} else if m.pressDuration > m.RepeatDelay {
					if (m.pressDuration-m.RepeatDelay)%m.RepeatInterval == 0 {
						m.OnRepeat()
					}
				}
			}

			x, y := m.touch.Pos()

			if m.OnDrag != nil {
				// ドラッグ動作：範囲外に出ても継続し、移動先を通知する
				m.OnDrag(x, y)
			} else {
				// クリック動作：範囲外に出たらキャンセルする
				cb := m.Widget.GetControlBase()
				gx, gy := cb.GetGlobalPos()
				// 範囲外に移動した判定
				if x < gx || gx+cb.Width <= x || y < gy || gy+cb.Height <= y {
					m.touch = nil
				}
			}
		}
	}

	cb := m.Widget.GetControlBase()
	gx, gy := cb.GetGlobalPos()
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

// updateFunctionは毎フレームの更新処理を行います。
// クリック完了、ドラッグ中の移動、ドラッグ終了の判定を行います。
func (m *MouseInteraction) updateFunction() {
}

// 現在のマウス/タッチ位置を取得する
func (m *MouseInteraction) GetPosition() (int, int, bool) {
	if m.touch != nil {
		x, y := m.touch.Pos()
		return x, y, true
	}
	return 0, 0, false
}
