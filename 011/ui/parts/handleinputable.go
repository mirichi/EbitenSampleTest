package parts

import "MyProject/ui/input"

// HandleInputableは独自の入力ロジックを持つ機能
type HandleInputable struct {
	Widget        Control
	OnHandleInput func(t input.Touch) bool
}

func NewHandleInputable(c Control) *HandleInputable {
	h := &HandleInputable{}
	h.InitHandleInputable(c)
	return h
}

func (h *HandleInputable) InitHandleInputable(c Control) {
	h.Widget = c

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
