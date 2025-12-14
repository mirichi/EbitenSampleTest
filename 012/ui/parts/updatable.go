package parts

// Updatableは独自の更新ロジックを追加したいときに使う
// 任意のタイミングで呼ばれるOnUpdateを定義できる
type Updatable struct {
	Control  Control
	OnUpdate func()
}

func NewUpdatable(c Control) *Updatable {
	u := &Updatable{}
	u.InitUpdatable(c)
	return u
}

func (u *Updatable) InitUpdatable(c Control) {
	u.Control = c

	// コントロールのUpdate時に呼ばれる関数を登録する
	c.GetControlBase().AddUpdateFunction(u.updateFunction)
}

// コントロールのUpdate時に呼ばれるUpdateFunction
func (u *Updatable) updateFunction() {
	if u.OnUpdate != nil {
		u.OnUpdate()
	}
}
