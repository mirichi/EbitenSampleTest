package main

import (
	"math"

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
	// 32x32 Blue Rect
	img := GetResourceManager().GetImage(ImgPlayer)

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
		p.X -= p.Speed
	}
	if input.IsActionPressed(input.ActionMoveRight) {
		p.X += p.Speed
	}
	if input.IsActionPressed(input.ActionMoveUp) {
		p.Y -= p.Speed
	}
	if input.IsActionPressed(input.ActionMoveDown) {
		p.Y += p.Speed
	}

	// Touch/Mouse Move
	t := input.GetPointer()
	if t != nil && t.IsPressed() {
		tx, ty := t.Pos()
		dx := tx - p.X
		dy := ty - p.Y
		dist := math.Hypot(dx, dy)

		if dist > p.Speed {
			p.X += (dx / dist) * p.Speed
			p.Y += (dy / dist) * p.Speed
		} else {
			p.X = tx
			p.Y = ty
		}
	}

	// Clamp to Screen
	w := float64(p.Image.Bounds().Dx())
	h := float64(p.Image.Bounds().Dy())

	if p.X < w/2 {
		p.X = w / 2
	}
	if p.X > ScreenWidth-w/2 {
		p.X = ScreenWidth - w/2
	}
	if p.Y < h/2 {
		p.Y = h / 2
	}
	if p.Y > ScreenHeight-h/2 {
		p.Y = ScreenHeight - h/2
	}

	p.Sprite.Update()
}

func (p *Player) IsDead() bool {
	return false // Player handles game over separately
}
