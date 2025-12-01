package parts

// Updatableは独自の更新ロジックを持つ機能
type Updatable struct {
	Control  Control
	OnUpdate func()
}

// Updatable生成
func NewUpdatable(c Control, f func()) *Updatable {
	u := &Updatable{
		Control:  c,
		OnUpdate: f,
	}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(u.updateFunction)

	return u
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (u *Updatable) updateFunction() {
	u.OnUpdate()
}
