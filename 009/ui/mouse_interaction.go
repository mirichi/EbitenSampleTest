package ui

// MouseInteractionはマウス/タッチ操作（クリック、ドラッグ、プレス）を扱う機能
type MouseInteraction struct {
	Control   Control
	touchInfo TouchInfo

	OnPress     func()
	OnClick     func()
	OnDragStart func(x, y int)
	OnDrag      func(dx, dy int)
	OnDragEnd   func()
}

// MouseInteraction生成
func NewMouseInteraction(c Control) *MouseInteraction {
	m := &MouseInteraction{
		Control: c,
	}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(m.updateFunction)

	return m
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (m *MouseInteraction) updateFunction(t TouchInfo) UpdateResult {
	var result UpdateResult = NotConsumed

	if m.touchInfo == nil { // 押されていないとき
		// tがnilの場合は、手前のコントロールで処理済みなので何もしない
		if t != nil {
			cb := m.Control.GetControlBase()
			ox, oy := cb.GetOwnerPos()
			x, y := t.Pos()
			// 座標判定
			if ox+cb.X <= x && x < ox+cb.X+cb.Width && oy+cb.Y <= y && y < oy+cb.Y+cb.Height {
				// 範囲内にある場合はイベントを消費する
				result = Consumed

				if t.IsJustPressed() { // 今回押された
					m.touchInfo = t
					if m.OnPress != nil {
						m.OnPress()
					}
					if m.OnDragStart != nil {
						m.OnDragStart(x, y)
					}
				}
			}
		}
	} else { // 押されているとき
		if m.touchInfo.IsJustReleased() { // 今回離された
			m.touchInfo = nil
			if m.OnClick != nil {
				m.OnClick()
			}
			if m.OnDragEnd != nil {
				m.OnDragEnd()
			}
		} else { // 継続して押されている
			x, y := m.touchInfo.Pos()

			if m.OnDrag != nil {
				// ドラッグ動作：範囲外に出ても継続し、移動量を通知する
				bx, by := m.touchInfo.OldPos()
				m.OnDrag(x-bx, y-by)
				result = Consumed
			} else {
				// クリック動作：範囲外に出たらキャンセルする
				cb := m.Control.GetControlBase()
				ox, oy := cb.GetOwnerPos()
				// 範囲外に移動した判定
				if x < ox+cb.X || ox+cb.X+cb.Width <= x || y < oy+cb.Y || oy+cb.Y+cb.Height <= y {
					m.touchInfo = nil
				} else {
					result = Consumed
				}
			}
		}
	}

	return result
}
