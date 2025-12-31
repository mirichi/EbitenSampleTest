package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/input"
	"MyProject/parts"
	"MyProject/ui"
	"MyProject/uiparts"
)

const (
	BoardWidth  = 10
	BoardHeight = 10
	BombCount   = 10
	CellSize    = 30
)

type Game struct {
	mainscreen *ui.MainScreen
	logic      *MinesweeperLogic
	timer      *parts.Timer
	timeCount  int
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
	ebiten.SetWindowSize(BoardWidth*CellSize+100, BoardHeight*CellSize+150)
	ebiten.SetWindowTitle("Minesweeper Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	logic := NewMinesweeperLogic(BoardWidth, BoardHeight, BombCount)

	// メインウィンドウ
	win := ui.NewWindow(20, 20, BoardWidth*CellSize+40, BoardHeight*CellSize+80, "Minesweeper")
	win.ClientArea.AutoLayout = uiparts.FlexLayoutV(uiparts.FlexStart, uiparts.FlexStretch, 10)

	// ステータスエリア
	statusPanel := ui.NewGroupingControl(0, 0, 0, 40)
	statusPanel.AutoLayout = uiparts.FlexLayoutH(uiparts.FlexSpaceAround, uiparts.FlexCenter, 0)
	win.AddChild(statusPanel)

	bombLabel := ui.NewLabel(0, 0, 60, 30, "Bombs: 10", 14)
	statusPanel.AddChild(bombLabel)

	resetBtn := ui.NewButton(0, 0, 60, 30, "Reset", 14)
	statusPanel.AddChild(resetBtn)

	timerLabel := ui.NewLabel(0, 0, 60, 30, "Time: 0", 14)
	statusPanel.AddChild(timerLabel)

	// 盤面エリア
	boardPanel := ui.NewGroupingControl(0, 0, BoardWidth*CellSize, BoardHeight*CellSize)
	win.AddChild(boardPanel)

	cells := make([][]*Cell, BoardHeight)
	for y := 0; y < BoardHeight; y++ {
		cells[y] = make([]*Cell, BoardWidth)
		for x := 0; x < BoardWidth; x++ {
			c := NewCell(float64(x)*CellSize, float64(y)*CellSize, CellSize, CellSize, x, y, logic)
			boardPanel.AddChild(c)
			cells[y][x] = c
		}
	}

	game := &Game{
		logic: logic,
	}

	// タイマー設定
	game.timer = parts.NewTimer(win, 60, true) // 60fps = 1sec
	game.timer.OnTimer = func() {
		if !logic.GameOver && !logic.GameWon && game.timeCount < 999 {
			game.timeCount++
			timerLabel.Text = fmt.Sprintf("Time: %d", game.timeCount)
		}
	}

	resetFunc := func() {
		logic.Reset()
		game.timeCount = 0
		timerLabel.Text = "Time: 0"
		game.timer.Start()
		bombLabel.Text = fmt.Sprintf("Bombs: %d", logic.BombCount)
	}

	resetBtn.OnClick = resetFunc

	// 更新時の状態監視
	win.AddAfterUpdateFunction(func() {
		bombCount := logic.BombCount - logic.GetFlagCount()
		bombLabel.Text = fmt.Sprintf("Bombs: %d", bombCount)

		if logic.GameOver {
			bombLabel.Text = "LOSE!"
			game.timer.Stop()
		} else if logic.GameWon {
			bombLabel.Text = "WIN!"
			game.timer.Stop()
		}
	})

	ms := ui.NewMainScreen()
	ms.AddChild(win)
	game.mainscreen = ms

	game.timer.Start()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
