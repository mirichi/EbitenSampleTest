package main

import (
	"image/color"
	"math"
	"math/rand/v2"

	"myproject/control"
	"myproject/primitive"
	"myproject/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	objects    []primitive.Object
	controls   []ui.Control
	dragMap    map[ui.TouchInfo]primitive.Object
	dragObj    map[primitive.Object]struct{}
	menuscreen *control.MenuScreen
}

func newGame() *Game {
	var game = &Game{}
	game.Init()
	return game
}

func (g *Game) Init() {
	// 管理用マップ生成
	g.dragMap = map[ui.TouchInfo]primitive.Object{}
	g.dragObj = map[primitive.Object]struct{}{}

	// Rect生成
	r1 := primitive.NewStar(120, 420, 100, math.Pi/4)
	r2 := primitive.NewStar(320, 120, 60, math.Pi/5)
	r3 := primitive.NewRect(400, 300, 150, 200, math.Pi/3)
	c1 := primitive.NewHarfCircle(200, 300, 50)
	c2 := primitive.NewCircle(100, 100, 40)
	g.objects = append(g.objects, r1, r2, r3, c1, c2)

	// メニュー画面生成
	g.menuscreen = control.NewMenuScreen(0, 0, 640, 480, nil)

	// メニューパネル生成
	menu := control.NewMenu(100, 50, 440, 380, g.menuscreen)

	// メニュー画面に登録
	g.menuscreen.Controls = append(g.menuscreen.Controls, menu)

	// メニューパネル上のラベル生成
	l := control.NewLabel(0, 10, 440, 0, "Menu", 50, ui.AdjustCenter, menu)

	// メニューパネル上のボタン生成
	mc1 := control.NewButton(120, 120, 200, 50, "1個増やす", 28, ui.AdjustCenter, menu, func() {
		g.objects = append(g.objects, primitive.NewRect(rand.Float64()*640, rand.Float64()*480, rand.Float64()*80+80, rand.Float64()*80+80, rand.Float64()*math.Pi))
	})
	mc2 := control.NewButton(120, 190, 200, 50, "1個減らす", 28, ui.AdjustCenter, menu, func() {
		if len(g.objects) > 5 {
			g.objects = g.objects[:len(g.objects)-1]
		}
	})

	// メニューパネルに登録
	menu.Controls = append(menu.Controls, mc1, mc2, l)

	// メニューボタン
	mb := control.NewButton(20, 20, 50, 50, "三", 28, ui.AdjustCenter, nil, func() {
		g.menuscreen.Start()
	})

	// メニューボタンをトップレベルに登録
	g.controls = append(g.controls, mb)

	// メニュー画面をトップレベルに登録
	g.controls = append(g.controls, g.menuscreen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, o := range g.objects {
		o.Draw(screen)
	}
	for _, o := range g.controls {
		o.Draw(screen)
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	ui.Input_Update()

	if !g.menuscreen.Running {
		// タッチ継続、終了の処理
		for tinfo, obj := range g.dragMap {
			// 押されている場合は移動処理
			if tinfo.IsPressed() {
				oldX, oldY := tinfo.OldPos()
				x, y := tinfo.Pos()
				obj.Move(float64(oldX), float64(oldY), float64(x), float64(y))
			} else {
				// 押されていない場合はマップから削除
				delete(g.dragMap, tinfo)
				delete(g.dragObj, obj)
			}
		}

		// タッチ開始の処理
		for _, tinfo := range ui.AllTouches() {
			// 今回押されたタッチ
			if tinfo.IsJustPressed() {
				x, y := tinfo.Pos()
				// 押されたRectを探す
				for _, obj := range g.objects {

					_, found := g.dragObj[obj]
					if !found && obj.CheckPoint(float64(x), float64(y)) {
						// ドラッグ中情報を保存
						g.dragMap[tinfo] = obj
						g.dragObj[obj] = struct{}{}
						break
					}
				}
			}
		}
	}

	// UI入力判定
	// もっと簡略化したい
	t := ui.FirstTouch()
	if t != nil && t.IsJustPressed() {
		for i := len(g.controls) - 1; i >= 0; i-- {
			if g.controls[i].ProcessTouch(t) {
				break
			}
		}
	}

	// ui処理
	for _, c := range g.controls {
		c.Update()
	}

	// 衝突判定用情報更新
	for _, r := range g.objects {
		r.Update()
	}

	// 衝突判定
	for i := 0; i < len(g.objects)-1; i++ {
		for j := i + 1; j < len(g.objects); j++ {
			r1 := g.objects[i]
			r2 := g.objects[j]
			if r1.TestCollinsion(r2) {
				r1.SetFillColor(color.RGBA{0xff, 0xff, 0x00, 0xff})
				r2.SetFillColor(color.RGBA{0xff, 0xff, 0x00, 0xff})
			}
		}
	}

	return nil
}

func main() {
	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
