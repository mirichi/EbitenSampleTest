package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/input"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type Game struct {
	Scene *GameScene
}

func NewGame() *Game {
	g := &Game{}
	g.Init()
	return g
}

func (g *Game) Init() {
	g.Scene = NewGameScene()
}

func (g *Game) Update() error {
	input.Update()
	g.Scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Shooting Game Sample")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
