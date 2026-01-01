package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/input"
	"MyProject/ui"
	"MyProject/uiparts"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type Game struct {
	Scene         *GameScene
	MainScreen    *ui.MainScreen
	ScoreLabel    *ui.Label
	BossHPLabel   *ui.Label
	GameOverLabel *ui.Label
}

func NewGame() *Game {
	g := &Game{}
	g.Init()
	return g
}

func (g *Game) Init() {
	g.Scene = NewGameScene()

	// UI初期化
	g.MainScreen = ui.NewMainScreen()

	// ステータス表示用のパネル
	// ゲーム画面の左上に配置
	hudPanel := ui.NewPanel(10, 10, 150, 50)
	hudPanel.BackgroundColor = color.RGBA{25, 25, 25, 255}
	// レイアウト設定（縦並び、左寄せ、幅ストレッチ、隙間5px）
	hudPanel.AutoLayout = uiparts.FlexLayoutV(uiparts.FlexStart, uiparts.FlexStretch, 5)

	// スコアラベル
	g.ScoreLabel = ui.NewLabel(0, 0, 0, 20, "Score: 0", 16)
	hudPanel.AddChild(g.ScoreLabel)

	// ボスHPラベル
	g.BossHPLabel = ui.NewLabel(0, 0, 0, 20, "Boss HP: -", 16)
	hudPanel.AddChild(g.BossHPLabel)

	g.MainScreen.AddChild(hudPanel)

	// リセットボタンなどを置くウィンドウ
	menuWin := ui.NewPanel(ScreenWidth-120, 10, 110, 40)
	menuWin.BackgroundColor = color.RGBA{25, 25, 25, 255}
	menuWin.AutoLayout = uiparts.FlexLayoutV(uiparts.FlexStart, uiparts.FlexStretch, 5)

	resetBtn := ui.NewButton(0, 0, 0, 30, "Reset", 16)
	resetBtn.OnPress = func() {
		g.Scene.InitGameScene()
	}
	menuWin.AddChild(resetBtn)

	g.MainScreen.AddChild(menuWin)

	// ゲームオーバーラベル（初期状態では親なし）
	g.GameOverLabel = ui.NewLabel(0, 200, 640, 50, "GAME OVER\n\rPress R to Reset", 48)
	// 中央揃えにするために、Label側の設定か、レイアウトを使うかだが、
	// NewLabelのデフォルトは左揃え。中央揃えにするにはLabelのプロパティをいじるか、
	// 単に位置調整で対応する。ここでは位置調整とテキストアラインメント変更。
	g.GameOverLabel.AlignX = uiparts.AlignCenter
}

func (g *Game) Update() error {
	input.Update()

	// UI更新
	g.MainScreen.HandleInput(input.GetPointer())
	g.MainScreen.Update()

	// UIがマウスをキャプチャしていない場合のみゲームを更新する... としたいところだが
	// 今のところクリック判定くらいしかないので、とりあえず常に更新する
	g.Scene.Update()

	// ゲームの状態をUIに反映
	g.ScoreLabel.Text = "Score: " + strconv.Itoa(g.Scene.Score)

	if g.Scene.BossSpawned && g.Scene.Boss != nil {
		g.BossHPLabel.Text = "HP: " + strconv.Itoa(g.Scene.Boss.HP)
	} else {
		g.BossHPLabel.Text = "HP: -"
	}

	// ゲームオーバー表示制御
	if g.Scene.GameOver {
		if g.GameOverLabel.GetEntityBase().Parent == nil {
			g.MainScreen.AddChild(g.GameOverLabel)
		}
	} else {
		if g.GameOverLabel.GetEntityBase().Parent != nil {
			g.MainScreen.RemoveChild(g.GameOverLabel)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
	g.MainScreen.Draw(screen)
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
