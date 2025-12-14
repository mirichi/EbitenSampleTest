package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/ui"
	"MyProject/ui/input"
	"MyProject/ui/parts"
)

type Game struct {
	mainscreen *ui.MainScreen
}

func (g *Game) Update() error {
	input.Update()
	g.mainscreen.HandleInput(input.GetPointer())
	g.mainscreen.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0x80, 0x80, 0xff})
	g.mainscreen.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Paint Tool Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// ウィンドウ生成
	win := ui.NewWindow(50, 20, 580, 350, "Paint Tool")
	win.ClientArea.AutoLayout = parts.AutoLayoutFitH

	// ColorPanel生成
	cp := NewColorPanel(0, 0, 50, 256)
	win.AddChild(cp)

	// スクロール可能なパネル生成
	sp := ui.NewScrollablePanel(0, 0, 0, 0, 15)
	sp.AutoResizable = true
	win.AddChild(sp)

	// キャンバス生成
	canvas := NewCanvas(0, 0, 480, 320)
	sp.AddChild(canvas)
	sp.SetMaxRange(&canvas.Width, &canvas.Height)

	// ToolBox生成
	tp := NewToolBox(0, 0, 50, 256)
	win.AddChild(tp)

	ms := ui.NewMainScreen()
	ms.AddChild(win)

	// ColorPanelの選択色でCanvasの色を変更する
	cp.OnSelect = func(c color.Color) {
		canvas.Color = c
	}

	// ToolBoxの選択サイズでCanvasの線幅を変更する
	tp.OnSelect = func(s int) {
		canvas.LineWidth = s
	}

	// 右クリックメニュー
	menu := ui.NewMenuBuilder().
		AddItem("Clear", func() { canvas.Clear() }).
		AddSeparator().
		AddSubMemu("Color", ui.NewMenuBuilder().
			AddItem("Red", func() { cp.Buttons[1].OnPress() }).
			AddItem("Green", func() { cp.Buttons[3].OnPress() }).
			AddItem("Blue", func() { cp.Buttons[5].OnPress() }).
			Build()).
		AddSubMemu("Size", ui.NewMenuBuilder().
			AddItem("2", func() { tp.Buttons[0].OnPress() }).
			AddItem("6", func() { tp.Buttons[2].OnPress() }).
			AddItem("10", func() { tp.Buttons[4].OnPress() }).
			Build()).
		Build()

	// menu := []*ui.MenuItem{
	// 	{Text: "Clear", Action: func() { canvas.Clear() }},
	// 	{Text: "Color", SubMenu: []*ui.MenuItem{
	// 		{Text: "Red", Action: func() { cp.Buttons[1].OnPress() }},
	// 		{Text: "Green", Action: func() { cp.Buttons[3].OnPress() }},
	// 		{Text: "Blue", Action: func() { cp.Buttons[5].OnPress() }},
	// 	}},
	// 	{Text: "Size", SubMenu: []*ui.MenuItem{
	// 		{Text: "2", Action: func() { tp.Buttons[0].OnPress() }},
	// 		{Text: "6", Action: func() { tp.Buttons[2].OnPress() }},
	// 		{Text: "10", Action: func() { tp.Buttons[4].OnPress() }},
	// 	}},
	// }

	canvas.OnRightRelease = func() {
		x, y := input.GetPointer().Pos()
		ms.PopupManager.ShowMenu(x, y, menu)
	}

	game := &Game{
		mainscreen: ms,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
