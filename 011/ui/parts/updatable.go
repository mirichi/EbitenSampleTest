package parts

// Updatableは独自の更新ロジックを持つ機能
type Updatable struct {
	Widget   Control
	OnUpdate func()
}

func NewUpdatable(c Control) *Updatable {
	u := &Updatable{}
	u.InitUpdatable(c)
	return u
}

func (u *Updatable) InitUpdatable(c Control) {
	u.Widget = c

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(u.updateFunction)
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (u *Updatable) updateFunction() {
	if u.OnUpdate != nil {
		u.OnUpdate()
	}
}
