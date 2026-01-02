package main

import (
	"MyProject/objects"
	"math"
)

type Bullet struct {
	*objects.Sprite
	*objects.CircleCollider
	vx, vy     float64
	isDeadFlag bool
}

func NewBullet(x, y float64) *Bullet {
	// Standard straight up bullet
	return NewBulletVelocity(x, y, 0, -8.0)
}

func NewBulletVelocity(x, y, vx, vy float64) *Bullet {
	// 8x8 Yellow Circle-ish
	// 8x8 Yellow Circle-ish
	img := GetResourceManager().GetImage(ImgBulletPlayer)

	s := objects.NewSprite(x, y, img)

	b := &Bullet{
		Sprite:     s,
		vx:         vx,
		vy:         vy,
		isDeadFlag: false,
	}

	// Circle Collider
	b.CircleCollider = objects.NewCircleCollider(s, objects.Vector2{X: 6, Y: 6}, 6)

	return b
}

func NewBulletAngled(x, y, angle, speed float64) *Bullet {
	vx := math.Cos(angle) * speed
	vy := math.Sin(angle) * speed
	return NewBulletVelocity(x, y, vx, vy)
}

func (b *Bullet) Update() {
	if b.isDeadFlag {
		return
	}

	b.Sprite.X += b.vx
	b.Sprite.Y += b.vy

	// Remove if out of screen (large margin to account for world scroll)
	if b.Sprite.Y < -50 || b.Sprite.X < -200 || b.Sprite.X > ScreenWidth+200 {
		b.isDeadFlag = true
	}

	b.Sprite.Update()
}

func (b *Bullet) IsDead() bool {
	return b.isDeadFlag
}

func (b *Bullet) MarkDead() {
	b.isDeadFlag = true
}

// Draw is promoted from Sprite
