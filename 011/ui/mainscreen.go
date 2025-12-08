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
	GroupingWidget
	parts.Updatable
	parts.Drawable
}

// MainScreen生成
func NewMainScreen() *MainScreen {
	ms := &MainScreen{}
	ms.InitUpdatable(ms)
	ms.InitGroupingWidget(0, 0, 0, 0)
	ms.InitDrawable(ms)
	ms.ClippingFlag = false
	ms.OrderChange = true

	PopupWidgets = NewPopupContainer()

	ms.OnUpdate = func() {
		PopupWidgets.Update()
		parts.FinalizeCursor()
	}

	ms.OnDraw = func(screen *ebiten.Image) {
		PopupWidgets.Draw(screen)
	}

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

	r := PopupWidgets.HandleInput(t)
	if !r {
		r = ms.GroupingWidget.HandleInput(t)
	}

	if !parts.FlameFocus && t != nil && t.IsJustPressed() && parts.FocusedWidget != nil {
		parts.FocusedWidget.Blur()
	}

	return r
}
