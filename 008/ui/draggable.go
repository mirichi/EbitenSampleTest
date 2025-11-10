package ui

// Draggableはドラッグ操作に反応する機能
type Draggable struct {
	control *ControlBase
	TouchInfo
	dragFunc func(dx, dy int)
}

// Draggable生成
func NewDraggable(c *ControlBase, f func(int, int)) *Draggable {
	r := &Draggable{
		control:  c,
		dragFunc: f,
	}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.AddUpdateFunction(r.UpdateFunction)

	return r
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (c *Draggable) UpdateFunction(t TouchInfo) {
	if c.TouchInfo == nil { // 押されていないとき
		if t != nil && t.IsJustPressed() { // 今回押された
			ox, oy := c.control.GetOwnerPos()
			x, y := t.Pos()
			// 座標判定
			if ox+c.control.X <= x && x < ox+c.control.X+c.control.Width && oy+c.control.Y <= y && y < oy+c.control.Y+c.control.Height {
				c.TouchInfo = t
			}
		}
	} else { // 押されているとき
		if c.TouchInfo.IsJustReleased() { // 今回離された
			c.TouchInfo = nil
		} else { // 継続して押されている
			x, y := t.Pos()
			bx, by := t.OldPos()
			c.dragFunc(x-bx, y-by)
		}
	}
}
