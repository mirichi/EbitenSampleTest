package parts

import "MyProject/input"

// HandleInputableは独自の入力ロジックを追加したいときに使う
// 任意のタイミングで呼ばれるOnHandleInputを定義できる
type HandleInputable struct {
	Control       Control
	OnHandleInput func(t input.Touch) bool
}

func NewHandleInputable(c Control) *HandleInputable {
	h := &HandleInputable{}
	h.InitHandleInputable(c)
	return h
}

func (h *HandleInputable) InitHandleInputable(c Control) {
	h.Control = c

	// コントロールのHandleInput時に呼ばれる関数を登録する
	c.GetControlBase().AddHandleInputFunction(h.handleInputFunction)
}

// コントロールのHandleInput時に呼ばれるHandleInputFunction
func (h *HandleInputable) handleInputFunction(t input.Touch) bool {
	if h.OnHandleInput != nil {
		return h.OnHandleInput(t)
	}
	return false
}
