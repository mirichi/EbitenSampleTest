package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/input"
	"MyProject/objects"
	"MyProject/ui"
	"MyProject/uiparts"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type Game struct {
	SceneManager   *objects.SceneManager
	MainScreen     *ui.MainScreen
	ScoreLabel     *ui.Label
	BossHPBar      *ui.ProgressBar
	GameOverLabel  *ui.Label
	GameOverLabel2 *ui.Label
}

func NewGame() *Game {
	g := &Game{}
	g.Init()
	return g
}

func (g *Game) Init() {
	// シーンマネージャー初期化
	// 最初はGameScene
	g.SceneManager = objects.NewSceneManager(NewGameScene())

	// UI初期化
	g.MainScreen = ui.NewMainScreen()

	// ステータス表示用のパネル (Score)
	// 画面左上に配置
	hudPanel := ui.NewPanel(10, 10, 150, 50)
	hudPanel.BackgroundColor = color.RGBA{25, 25, 25, 255}
	// レイアウト設定（縦並び、左寄せ、幅ストレッチ、隙間5px）
	hudPanel.AutoLayout = uiparts.FlexLayoutV(uiparts.FlexStart, uiparts.FlexStretch, 5)

	// スコアラベル
	g.ScoreLabel = ui.NewLabel(0, 0, 0, 20, "Score: 0", 16)
	hudPanel.AddChild(g.ScoreLabel)

	g.MainScreen.AddChild(hudPanel)

	// Boss HP Bar (Top Center, between Score and Reset)
	// Score ends at 160, Reset starts at 520. Available: 360px.
	// Bar Width: 300px. Center at 320. Start X = 170.
	g.BossHPBar = ui.NewProgressBar(0, 50, 300, 15)
	g.BossHPBar.X = 170
	g.BossHPBar.Y = 10
	g.BossHPBar.FillColor = color.RGBA{200, 50, 50, 255} // Red
	g.BossHPBar.Visible = false
	g.MainScreen.AddChild(g.BossHPBar)

	// リセットボタンなどを置くウィンドウ
	menuWin := ui.NewPanel(ScreenWidth-120, 10, 110, 40)
	menuWin.BackgroundColor = color.RGBA{25, 25, 25, 255}
	menuWin.AutoLayout = uiparts.FlexLayoutV(uiparts.FlexStart, uiparts.FlexStretch, 5)

	resetBtn := ui.NewButton(0, 0, 0, 30, "Reset", 16)
	resetBtn.OnPress = func() {
		if gs, ok := g.SceneManager.GetCurrentScene().(*GameScene); ok {
			gs.InitGameScene()
		}
	}
	menuWin.AddChild(resetBtn)

	g.MainScreen.AddChild(menuWin)

	// ゲームオーバーラベル（初期状態では親なし）
	g.GameOverLabel = ui.NewLabel(0, 200, 640, 50, "GAME OVER", 48)
	g.GameOverLabel2 = ui.NewLabel(0, 250, 640, 50, "Press R to Reset", 48)
	// 中央揃えにするために、Label側の設定か、レイアウトを使うかだが、
	// NewLabelのデフォルトは左揃え。中央揃えにするにはLabelのプロパティをいじるか、
	// 単に位置調整で対応する。ここでは位置調整とテキストアラインメント変更。
	g.GameOverLabel.AlignX = uiparts.AlignCenter
	g.GameOverLabel2.AlignX = uiparts.AlignCenter
}

func (g *Game) Update() error {
	input.Update()

	// UI更新
	g.MainScreen.HandleInput(input.GetPointer())
	g.MainScreen.Update()

	// シーン更新
	if err := g.SceneManager.Update(); err != nil {
		return err
	}

	// ゲームの状態をUIに反映
	// 現在のシーンがGameSceneの場合のみ反映
	if gs, ok := g.SceneManager.GetCurrentScene().(*GameScene); ok {
		g.ScoreLabel.Text = "Score: " + strconv.Itoa(gs.Score)

		if gs.BossSpawned && gs.Boss != nil {
			g.BossHPBar.Visible = true
			g.BossHPBar.SetRange(0, float64(gs.Boss.GetMaxHP()))
			g.BossHPBar.SetValue(float64(gs.Boss.GetHP()))
		} else {
			g.BossHPBar.Visible = false
		}

		// ゲームオーバー表示制御
		if gs.GameOver {
			if g.GameOverLabel.GetEntityBase().Parent == nil {
				g.MainScreen.AddChild(g.GameOverLabel)
				g.MainScreen.AddChild(g.GameOverLabel2)
			}
		} else {
			if g.GameOverLabel.GetEntityBase().Parent != nil {
				g.MainScreen.RemoveChild(g.GameOverLabel)
				g.MainScreen.RemoveChild(g.GameOverLabel2)
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.SceneManager.Draw(screen)
	g.MainScreen.Draw(screen)

	// Debug Info
	// ebitenutil.DebugPrint(screen, input.GetGamepadDebugInfo())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Shooting Game Sample with GUI")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
