package ui

import (
	"MyProject/ui/input"
	"MyProject/ui/parts"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

// MainScreenはオブジェクトツリーの最上位になる構造体
// 全てのWidgetをこのMainScreenの子要素として登録することで、
// MainScreenのHandleInput/Update/Drawを呼び出すだけでUI全体の更新と描画が行われる
type MainScreen struct {
	BlankWidget
}

// MainScreen生成
func NewMainScreen() *MainScreen {
	ms := &MainScreen{}
	ms.InitBlankWidget(nil, 0, 0, 0, 0)
	ms.ClippingFlag = false
	ms.OrderChange = true

	return ms
}

// HandleInputは入力イベントを処理する
// ウィンドウサイズに合わせてサイズを更新し、フォーカス制御を行う
func (ms *MainScreen) HandleInput(t input.Touch) bool {
	// ブラウザではウィンドウサイズが取得できないので固定値にする
	if runtime.GOOS == "js" {
		ms.Width = 640
		ms.Height = 480
	} else {
		ms.Width, ms.Height = ebiten.WindowSize()
	}
	parts.FlameFocus = false
	r := ms.WidgetBase.HandleInput(t)
	if !parts.FlameFocus && t != nil && t.IsJustPressed() && parts.FocusedWidget != nil {
		parts.FocusedWidget.Blur()
	}
	return r
}

func (ms *MainScreen) Update() {
	ms.WidgetBase.Update()
	parts.FinalizeCursor()
}
