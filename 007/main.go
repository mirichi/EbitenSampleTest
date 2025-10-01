package main

import (
	"myproject/control"
	"myproject/primitive"
	"myproject/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	objects  []primitive.Object
	controls []ui.Control
	dragMap  map[ui.TouchInfo]primitive.Object
	dragObj  map[primitive.Object]struct{}
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
	c1 := primitive.NewSimpleCircle(80, 300, 10)
	c2 := primitive.NewSimpleCircle(250, 50, 10)
	c3 := primitive.NewSimpleCircle(500, 400, 10)
	c4 := primitive.NewSimpleCircle(600, 100, 10)
	g.objects = append(g.objects, c1, c2, c3, c4)

	s1 := control.NewSlider(270, 440, 100, 40, "50", 24, ui.AdjustCenter, nil, func() {
		g.controls[0].(*control.Slider).Slide() // 自分自身をアクセスする手段が無く苦肉の策
	})
	g.controls = append(g.controls, s1)

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

	g.draw_bezier(screen)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	ui.Input_Update()

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

	return nil
}

func main() {
	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
