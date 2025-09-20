package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type drawer interface {
	Draw(screen *ebiten.Image)
}

type Game struct {
	objects    []drawer
	oldX, oldY int
	dragRect   *Rect
	id         ebiten.TouchID
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
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	x, y := ebiten.CursorPosition()
	touchIDs := inpututil.AppendJustPressedTouchIDs(nil)

	// マウスでRectオブジェクトをドラッグする
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, o := range g.objects {
			if o.(*Rect).CheckPoint(float64(x), float64(y)) {
				g.oldX, g.oldY = x, y
				g.dragRect = o.(*Rect)
				break
			}
		}
	} else if len(touchIDs) > 0 {
		g.id = touchIDs[0]
		tx, ty := ebiten.TouchPosition(g.id)
		for _, o := range g.objects {
			if o.(*Rect).CheckPoint(float64(tx), float64(ty)) {
				g.oldX, g.oldY = tx, ty
				g.dragRect = o.(*Rect)
				break
			}
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.dragRect != nil {
			g.dragRect.Move(float64(g.oldX), float64(g.oldY), float64(x), float64(y))
			g.oldX, g.oldY = x, y
		}
	} else if len(touchIDs) > 0 && g.id == touchIDs[0] {
		tx, ty := ebiten.TouchPosition(g.id)
		if g.dragRect != nil {
			g.dragRect.Move(float64(g.oldX), float64(g.oldY), float64(tx), float64(ty))
			g.oldX, g.oldY = tx, ty
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.dragRect != nil {
			g.dragRect = nil
		}
	} else if g.id != -1 && inpututil.IsTouchJustReleased(g.id) {
		if g.dragRect != nil {
			g.dragRect = nil
			g.id = -1
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
