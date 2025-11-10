package ui

// MainScreenはコントロールツリーの最上位になる構造体
// 機能はGroupingだけ
// コントロールをこれに登録することで、MainScreenのUpdate/Drawをするだけでよくなる
type MainScreen struct {
	*ControlBase
	*Grouping
}

// MainScreen生成
func NewMainScreen() *MainScreen {
	cb := NewControlBase(0, 0, 0, 0)

	return &MainScreen{
		ControlBase: cb,
		Grouping:    NewGrouping(cb),
	}
}
