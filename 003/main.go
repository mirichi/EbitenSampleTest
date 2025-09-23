package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type drawer interface {
	Draw(screen *ebiten.Image)
}

type Game struct {
	objects  []drawer
	touchMap map[TouchInfo]drawer
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
	g.touchMap = make(map[TouchInfo]drawer)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, o := range g.objects {
		o.Draw(screen)
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	Input_Update()

	// タッチ継続、終了の処理
	for tinfo, obj := range g.touchMap {
		// 押されている場合は移動処理
		if tinfo.IsPressed() {
			rect := obj.(*Rect)
			oldX, oldY := tinfo.OldPos()
			x, y := tinfo.Pos()
			rect.Move(float64(oldX), float64(oldY), float64(x), float64(y))
		} else {
			// 押されていない場合はマップから削除
			delete(g.touchMap, tinfo)
		}
	}

	// タッチ開始の処理
	for _, tinfo := range AllTouches() {
		if tinfo.IsJustPressed() {
			x, y := tinfo.Pos()
			for _, o := range g.objects {
				rect := o.(*Rect)
				if rect.CheckPoint(float64(x), float64(y)) {
					g.touchMap[tinfo] = rect
					break
				}
			}
		}
	}

	// Rect衝突判定
	if g.objects[0].(*Rect).CheckBox(g.objects[1].(*Rect)) {
		g.objects[0].(*Rect).FillColor = color.RGBA{0xff, 0xff, 0x00, 0xff}
		g.objects[1].(*Rect).FillColor = color.RGBA{0xff, 0xff, 0x00, 0xff}
	} else {
		g.objects[0].(*Rect).FillColor = color.RGBA{0x00, 0xff, 0xff, 0xff}
		g.objects[1].(*Rect).FillColor = color.RGBA{0x00, 0xff, 0xff, 0xff}
	}

	return nil
}

func main() {
	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
