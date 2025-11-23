package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// MainScreenはコントロールツリーの最上位になる構造体
// コントロールをこれに登録することで、MainScreenのUpdate/Drawをするだけでよくなる
type MainScreen struct {
	*ControlBase
	*Grouping
}

// MainScreen生成
func NewMainScreen() *MainScreen {
	ms := &MainScreen{}
	ms.ControlBase = NewControlBase(ms, 0, 0, 0, 0)
	ms.Grouping = NewGrouping(ms)
	ms.Grouping.OrderChange = true

	return ms
}

func (ms *MainScreen) Update(t TouchInfo) UpdateResult {
	ms.ControlBase.Width, ms.ControlBase.Height = ebiten.WindowSize()
	FlameFocus = false
	r := ms.ControlBase.Update(t)
	if !FlameFocus && t != nil && t.IsJustPressed() && FocusedControl != nil {
		FocusedControl.Blur()
	}
	return r
}
