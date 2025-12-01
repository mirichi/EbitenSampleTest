package parts

import "MyProject/ui/input"

// HandleInputableは独自の入力ロジックを持つ機能
type HandleInputable struct {
	Control       Control
	OnHandleInput func(t input.Touch) bool
}

// HandleInputable生成
func NewHandleInputable(c Control, f func(t input.Touch) bool) *HandleInputable {
	u := &HandleInputable{
		Control:       c,
		OnHandleInput: f,
	}

	// コントロールのHandleInput時に呼ばれる関数を登録する
	c.GetControlBase().AddHandleInputFunction(u.handleInputFunction)

	return u
}

// コントロールのHandleInput時に呼ばれるHandleInputFunction
func (u *HandleInputable) handleInputFunction(t input.Touch) bool {
	return u.OnHandleInput(t)
}
