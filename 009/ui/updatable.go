package ui

// Updatableは独自の更新ロジックを持つ機能
type Updatable struct {
	Control Control
	f       func(t TouchInfo) UpdateResult
}

// Updatable生成
func NewUpdatable(c Control, f func(t TouchInfo) UpdateResult) *Updatable {
	u := &Updatable{
		Control: c,
		f:       f,
	}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(u.updateFunction)

	return u
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (u *Updatable) updateFunction(t TouchInfo) UpdateResult {
	return u.f(t)
}
