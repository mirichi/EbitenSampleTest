package parts

// Updatableは独自の更新ロジックを持つ機能
type Updatable struct {
	Widget   Widget
	OnUpdate func()
}

func NewUpdatable(c Widget) *Updatable {
	u := &Updatable{}
	u.InitUpdatable(c)
	return u
}

func (u *Updatable) InitUpdatable(c Widget) {
	u.Widget = c

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetWidgetBase().AddUpdateFunction(u.updateFunction)
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (u *Updatable) updateFunction() {
	if u.OnUpdate != nil {
		u.OnUpdate()
	}
}
