package main

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type drawer interface {
	Draw(screen *ebiten.Image)
}

type Game struct {
	objects  []drawer
	controls []control
	dragMap  map[TouchInfo]drawer
}

func newGame() *Game {
	var game = &Game{}
	game.Init()
	return game
}

func (g *Game) Init() {
	r1 := NewRect(220, 220, 100, 80)
	r2 := NewRect(400, 300, 150, 200)
	g.objects = append(g.objects, r1, r2)

	c1 := NewButton(10, 10, 200, 50, "1個増やす", func() {
		g.objects = append(g.objects, NewRect(rand.Float64()*640, rand.Float64()*480, 100, 100))
	})
	c2 := NewButton(10, 70, 200, 50, "1個増減らす", func() {
		if len(g.objects) > 2 {
			g.objects = g.objects[:len(g.objects)-1]
		}
	})
	g.controls = append(g.controls, c1, c2)

	g.dragMap = make(map[TouchInfo]drawer)
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

	Input_Update()

	// タッチ継続、終了の処理
	for tinfo, obj := range g.dragMap {
		// 押されている場合は移動処理
		if tinfo.IsPressed() {
			rect := obj.(*Rect)
			oldX, oldY := tinfo.OldPos()
			x, y := tinfo.Pos()
			rect.Move(float64(oldX), float64(oldY), float64(x), float64(y))
		} else {
			// 押されていない場合はマップから削除
			delete(g.dragMap, tinfo)
		}
	}

	// タッチ開始の処理
	for _, tinfo := range AllTouches() {
		// 今回押されたタッチ
		if tinfo.IsJustPressed() {
			x, y := tinfo.Pos()
			// 押されたRectを探す
			for _, o := range g.objects {
				rect := o.(*Rect)
				if rect.CheckPoint(float64(x), float64(y)) {
					// ドラッグ中情報を保存
					g.dragMap[tinfo] = rect
					break
				}
			}
		}
	}

	// UI入力判定
	at := AllTouches()
	if len(at) > 0 {
		t := at[0]
		if t.IsJustPressed() {
			for _, c := range g.controls {
				if c.CheckPoint(t.Pos()) {
					c.Press(t)
				}
			}
		}
	}

	// ui処理
	for _, c := range g.controls {
		c.Update()
	}

	// Rect衝突判定
	for _, r := range g.objects {
		r.(*Rect).FillColor = color.RGBA{0x00, 0xff, 0xff, 0xff}
	}

	for _, r1 := range g.objects {
		for _, r2 := range g.objects {
			if r1 != r2 {
				if r1.(*Rect).CheckBox(r2.(*Rect)) {
					r1.(*Rect).FillColor = color.RGBA{0xff, 0xff, 0x00, 0xff}
					r2.(*Rect).FillColor = color.RGBA{0xff, 0xff, 0x00, 0xff}
				}
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
