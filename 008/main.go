package main

import (
	"MyProject/ui"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	mainscreen *ui.MainScreen
)

type Game struct{}

func (g *Game) Update() error {
	// 入力更新
	ui.Input_Update()

	// コントロールのUpdate
	mainscreen.Update(ui.FirstTouch())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// コントロールの描画
	mainscreen.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	// MainScreen生成
	mainscreen = ui.NewMainScreen()

	// Window生成、MainScreenに登録
	window := ui.NewWindow(200, 100, 300, 300, "TestWindow")
	mainscreen.AddChild(window)

	// Button生成、Windowに登録
	var button1, button2 *ui.Button
	var count1, count2 int
	button1 = ui.NewButton(50, 50, 200, 50, "テスト", 30, func() { count1++; button1.Text = fmt.Sprintf("%d回おした", count1) })
	button2 = ui.NewButton(50, 150, 200, 50, "テスト", 30, func() { count2++; button2.Text = fmt.Sprintf("%d回おした", count2) })
	window.AddChild(button1)
	window.AddChild(button2)

	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
