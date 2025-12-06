package parts

// Updatableは独自の更新ロジックを持つ機能
type Updatable struct {
	Widget   Widget
	OnUpdate func()
}

// Updatable生成
func NewUpdatable(c Widget, f func()) *Updatable {
	u := &Updatable{
		Widget:   c,
		OnUpdate: f,
	}

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetWidgetBase().AddUpdateFunction(u.updateFunction)

	return u
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (u *Updatable) updateFunction() {
	u.OnUpdate()
}
