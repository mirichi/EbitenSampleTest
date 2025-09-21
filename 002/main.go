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
	dragRect *Rect
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, o := range g.objects {
		o.Draw(screen)
	}
	Input_Draw(screen)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	Input_Update()

	// マウスでRectオブジェクトをドラッグする
	if IsButtonJustPressed() {
		x, y := CurrectPos()
		for _, o := range g.objects {
			if o.(*Rect).CheckPoint(float64(x), float64(y)) {
				g.dragRect = o.(*Rect)
				break
			}
		}
	} else if IsButtonPressed() {
		if g.dragRect != nil {
			x, y := CurrectPos()
			oldX, oldY := OldPos()
			g.dragRect.Move(float64(oldX), float64(oldY), float64(x), float64(y))
		}
	} else if IsButtonJustReleased() {
		if g.dragRect != nil {
			g.dragRect = nil
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
