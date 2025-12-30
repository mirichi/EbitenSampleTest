package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
)

type Player struct {
	*objects.Sprite
	*objects.PolygonCollider
	Speed float64
}

func NewPlayer() *Player {
	// 32x32 Blue Rect
	img := ebiten.NewImage(48, 48)
	img.Fill(color.RGBA{0, 0, 255, 255})

	// Start at bottom center
	s := objects.NewSprite(ScreenWidth/2, ScreenHeight-50, img)

	p := &Player{
		Sprite: s,
		Speed:  4.0,
	}

	// Collision
	p.PolygonCollider = objects.NewRectCollider(s)

	return p
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Sprite.X -= p.Speed
		if p.Sprite.X < float64(p.Sprite.Image.Bounds().Dx())/2 {
			p.Sprite.X = float64(p.Sprite.Image.Bounds().Dx()) / 2
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Sprite.X += p.Speed
		if p.Sprite.X > ScreenWidth-float64(p.Sprite.Image.Bounds().Dx())/2 {
			p.Sprite.X = ScreenWidth - float64(p.Sprite.Image.Bounds().Dx())/2
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Sprite.Y -= p.Speed
		if p.Sprite.Y < float64(p.Sprite.Image.Bounds().Dy())/2 {
			p.Sprite.Y = float64(p.Sprite.Image.Bounds().Dy()) / 2
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Sprite.Y += p.Speed
		if p.Sprite.Y > ScreenHeight-float64(p.Sprite.Image.Bounds().Dy())/2 {
			p.Sprite.Y = ScreenHeight - float64(p.Sprite.Image.Bounds().Dy())/2
		}
	}

	p.Sprite.Update()
}

// Draw is promoted from Sprite
// func (p *Player) Draw(screen *ebiten.Image) {
// 	p.Sprite.Draw(screen)
// }

func (p *Player) X() float64 { return p.Sprite.X }
func (p *Player) Y() float64 { return p.Sprite.Y }
