package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/input"

	"MyProject/objects"
)

type Player struct {
	*objects.Sprite
	objects.CollisionTester
	Speed      float64
	PowerLevel int
}

func NewPlayer() *Player {
	// 32x32 Blue Rect
	img := ebiten.NewImage(48, 48)
	img.Fill(color.RGBA{0, 0, 255, 255})

	// Start at bottom center
	s := objects.NewSprite(ScreenWidth/2, ScreenHeight-50, img)

	p := &Player{
		Sprite:     s,
		Speed:      4.0,
		PowerLevel: 1,
	}

	// Collision
	p.CollisionTester = objects.NewCircleCollider(s, objects.Vector2{X: 24, Y: 24}, 10)

	return p
}

func (p *Player) PowerUp() {
	p.PowerLevel++
	if p.PowerLevel > 4 {
		p.PowerLevel = 4
	}
}

func (p *Player) Update() {
	// Keyboard Move
	if input.IsActionPressed(input.ActionMoveLeft) {
		p.Sprite.X -= p.Speed
	}
	if input.IsActionPressed(input.ActionMoveRight) {
		p.Sprite.X += p.Speed
	}
	if input.IsActionPressed(input.ActionMoveUp) {
		p.Sprite.Y -= p.Speed
	}
	if input.IsActionPressed(input.ActionMoveDown) {
		p.Sprite.Y += p.Speed
	}

	// Touch/Mouse Move
	t := input.GetPointer()
	if t != nil && t.IsPressed() {
		tx, ty := t.Pos()
		dx := tx - p.Sprite.X
		dy := ty - p.Sprite.Y
		dist := math.Hypot(dx, dy)

		if dist > p.Speed {
			p.Sprite.X += (dx / dist) * p.Speed
			p.Sprite.Y += (dy / dist) * p.Speed
		} else {
			p.Sprite.X = tx
			p.Sprite.Y = ty
		}
	}

	// Clamp to Screen
	w := float64(p.Sprite.Image.Bounds().Dx())
	h := float64(p.Sprite.Image.Bounds().Dy())

	if p.Sprite.X < w/2 {
		p.Sprite.X = w / 2
	}
	if p.Sprite.X > ScreenWidth-w/2 {
		p.Sprite.X = ScreenWidth - w/2
	}
	if p.Sprite.Y < h/2 {
		p.Sprite.Y = h / 2
	}
	if p.Sprite.Y > ScreenHeight-h/2 {
		p.Sprite.Y = ScreenHeight - h/2
	}

	p.Sprite.Update()
}

// Draw is promoted from Sprite
// func (p *Player) Draw(screen *ebiten.Image) {
// 	p.Sprite.Draw(screen)
// }

func (p *Player) X() float64 { return p.Sprite.X }
func (p *Player) Y() float64 { return p.Sprite.Y }
