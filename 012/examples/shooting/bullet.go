package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
)

type Bullet struct {
	*objects.Sprite
	*objects.CircleCollider
	Speed  float64
	IsDead bool
}

func NewBullet(x, y float64) *Bullet {
	// 8x8 Yellow Circle-ish
	img := ebiten.NewImage(8, 8)
	img.Fill(color.RGBA{255, 255, 0, 255})

	s := objects.NewSprite(x, y, img)
	s.OriginX = 4
	s.OriginY = 4

	b := &Bullet{
		Sprite: s,
		Speed:  8.0,
		IsDead: false,
	}

	// Circle Collider
	b.CircleCollider = objects.NewCircleCollider(s, objects.Vector2{X: 4, Y: 4}, 4)

	return b
}

func (b *Bullet) Update() {
	if b.IsDead {
		return
	}

	b.Sprite.Y -= b.Speed

	// Remove if out of screen
	if b.Sprite.Y < -10 {
		b.IsDead = true
	}

	b.Sprite.Update()
}

// Draw is promoted from Sprite
