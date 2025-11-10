package ui

// Clickableはクリック操作に反応する機能
type Clickable struct {
	control *ControlBase
	TouchInfo
	tapFunc func()
}

// Cliclable生成
// fはクリック時に実行される関数
func NewClickable(c *ControlBase, f func()) *Clickable {
	r := &Clickable{
		control: c,
		tapFunc: f,
	}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.AddUpdateFunction(r.UpdateFunction)

	return r
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (c *Clickable) UpdateFunction(t TouchInfo) {
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
			c.tapFunc()
		} else { // 継続して押されている
			ox, oy := c.control.GetOwnerPos()
			x, y := t.Pos()
			// 範囲外に移動した判定
			if x < ox+c.control.X || ox+c.control.X+c.control.Width <= x || y < oy+c.control.Y || oy+c.control.Y+c.control.Height <= y {
				c.TouchInfo = nil
			}
		}
	}
}
