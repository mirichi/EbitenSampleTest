// mainパッケージはアプリケーションのエントリーポイント
package main

import (
	"MyProject/parts"
	"MyProject/ui"
	"MyProject/ui/input"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	mainscreen *ui.MainScreen
)

// GameはEbitengineのGameインターフェースを実装する構造体
type Game struct{}

// Updateは毎フレーム呼び出される更新処理
// 入力の更新、UIの更新、カーソル形状の適用を行う
func (g *Game) Update() error {
	// 入力更新
	input.Update()

	// コントロールの更新
	mainscreen.HandleInput(input.GetPointer())
	mainscreen.Update()

	return nil
}

// Drawは毎フレーム呼び出される描画処理
func (g *Game) Draw(screen *ebiten.Image) {
	// コントロールの描画
	mainscreen.Draw(screen)
}

// Layoutはウィンドウサイズが変更されたときに呼び出される
// ゲームの論理画面サイズを返す
func (a *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

// func (a *Game) LayoutF(outsideWidth, outsideHeight float64) (float64, float64) {
// 	s := deviceScaleFactor()
// 	return 640 * s, 480 * s
// }

// createWindowはテスト用のウィンドウを作成して返す
// ボタン、チェックボックス、テキストボックスを含む
func createWindow(x, y, w, h float64, str string) *ui.Window {
	// Window生成
	window := ui.NewWindow(x, y, w, h, str)

	// Button生成、Windowに登録
	var count1, count2 int
	button1 := ui.NewButton(50, 50, 200, 50, "テスト", 30)
	button1.OnClick = func() { count1++; button1.Text = fmt.Sprintf("%d回おした", count1) }
	button2 := ui.NewButton(50, 150, 200, 50, "テスト", 30)
	button2.OnClick = func() { count2++; button2.Text = fmt.Sprintf("%d回おした", count2) }
	window.AddChild(button1)
	window.AddChild(button2)

	// Checkbox生成、Windowに登録
	checkbox := ui.NewCheckbox(50, 250, 30, "OFF", 20, false)
	checkbox.OnCheckChanged = func(c bool) {
		if c {
			checkbox.Text = "ON"
		} else {
			checkbox.Text = "OFF"
		}
	}
	window.AddChild(checkbox)

	// TextBox生成、Windowに登録
	textbox := ui.NewTextBox(50, 300, 200, 30, "日本語入力", 20)
	window.AddChild(textbox)

	// ClientAreaのAutoLayoutを設定してオートレイアウト実行
	// window.ClientArea.AutoLayout = parts.AutoLayoutV
	window.ClientArea.AutoLayout = parts.FlexLayoutV(parts.FlexSpaceAround, parts.FlexCenter, 0)

	return window
}

func main() {
	// MainScreen生成
	mainscreen = ui.NewMainScreen()

	mainscreen.AddChild(createWindow(200, 100, 300, 300, "TestWindow1"))
	mainscreen.AddChild(createWindow(400, 300, 300, 300, "TestWindow2"))

	label := ui.NewLabel(50, 50, 200, 50, "Label", 30)
	mainscreen.AddChild(label)

	labelclicker := parts.NewMouseInteraction(label)
	labelclicker.OnClick = func() {
		label.Text = "Clicked"
	}

	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}

// func deviceScaleFactor() float64 {
// 	return ebiten.Monitor().DeviceScaleFactor()
// }
