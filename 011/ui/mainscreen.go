package ui

import (
	"MyProject/ui/input"
	"MyProject/ui/parts"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

// MainScreenはコントロールツリーの最上位になる構造体
// 全てのコントロールをこのMainScreenの子要素として登録することで、
// MainScreenのUpdate/Drawを呼び出すだけでUI全体の更新と描画が行われる
type MainScreen struct {
	parts.ControlBase
	parts.Grouping
}

// MainScreen生成
func NewMainScreen() *MainScreen {
	ms := &MainScreen{}
	ms.InitControlBase(ms, 0, 0, 0, 0)
	ms.InitGrouping(ms)
	ms.ClippingFlag = false
	ms.OrderChange = true

	return ms
}

// HandleInputは入力イベントを処理する
// ウィンドウサイズに合わせてサイズを更新し、フォーカス制御を行う
func (ms *MainScreen) HandleInput(t input.Touch) bool {
	// ブラウザではウィンドウサイズが取得できないので固定値にする
	if runtime.GOOS == "js" {
		ms.ControlBase.Width = 640
		ms.ControlBase.Height = 480
	} else {
		ms.ControlBase.Width, ms.ControlBase.Height = ebiten.WindowSize()
	}
	parts.FlameFocus = false
	r := ms.ControlBase.HandleInput(t)
	if !parts.FlameFocus && t != nil && t.IsJustPressed() && parts.FocusedControl != nil {
		parts.FocusedControl.Blur()
	}
	return r
}

func (ms *MainScreen) Update() {
	ms.ControlBase.Update()
	parts.FinalizeCursor()
}
