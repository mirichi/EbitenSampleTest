package ui

import (
	"MyProject/input"
	"MyProject/uiparts"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

// MainScreenはオブジェクトツリーの最上位になる構造体
// 全てのControlをこのMainScreenの子要素として登録することで、
// MainScreenのHandleInput/Update/Drawを呼び出すだけでUI全体の更新と描画が行われる
type MainScreen struct {
	GroupingControl

	PopupManager *PopupManager
}

// MainScreen生成
func NewMainScreen() *MainScreen {
	ms := &MainScreen{}
	ms.InitGroupingControl(0, 0, 0, 0)
	ms.OrderChange = true

	ms.PopupManager = NewPopupManager()

	ms.AddBeforeUpdateFunction(func() {
		ms.PopupManager.Update()
		ms.Layout()
		ms.PopupManager.Layout()
		uiparts.FinalizeCursor()
	})

	ms.AddDrawFunction(func(screen *ebiten.Image) {
		ms.PopupManager.Draw(screen)
	})

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
		w, h := ebiten.WindowSize()
		ms.Width, ms.Height = float64(w), float64(h)
	}

	uiparts.FlameFocus = false

	r := ms.PopupManager.HandleInput(t)
	if !r {
		r = ms.GroupingControl.HandleInput(t)
	}

	if !uiparts.FlameFocus && t != nil && t.IsJustPressed() && uiparts.FocusedControl != nil {
		uiparts.FocusedControl.Blur()
	}

	return r
}
