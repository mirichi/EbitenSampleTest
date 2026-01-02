package main

import (
	"math"

	"MyProject/objects"
)

type EnemyBullet struct {
	*objects.Sprite
	*objects.CircleCollider
	SpeedX     float64
	SpeedY     float64
	isDeadFlag bool
}

func NewEnemyBullet(x, y float64, angle float64) *EnemyBullet {
	// 12x12 Pink/Red Circle
	// 12x12 Pink/Red Circle
	img := GetResourceManager().GetImage(ImgBulletEnemy)

	s := objects.NewSprite(x, y, img)

	speed := 4.0
	vx := math.Cos(angle) * speed
	vy := math.Sin(angle) * speed

	b := &EnemyBullet{
		Sprite:     s,
		SpeedX:     vx,
		SpeedY:     vy,
		isDeadFlag: false,
	}

	b.CircleCollider = objects.NewCircleCollider(s, objects.Vector2{X: 6, Y: 6}, 6)

	return b
}

func (b *EnemyBullet) Update() {
	if b.isDeadFlag {
		return
	}

	b.Sprite.X += b.SpeedX
	b.Sprite.Y += b.SpeedY

	// Remove if out of screen (large margin)
	if b.Sprite.Y > 480+50 || b.Sprite.X < -200 || b.Sprite.X > ScreenWidth+200 || b.Sprite.Y < -100 {
		b.isDeadFlag = true
	}

	b.Sprite.Update()
}

func (b *EnemyBullet) IsDead() bool {
	return b.isDeadFlag
}

func (b *EnemyBullet) MarkDead() {
	b.isDeadFlag = true
}

// Draw is promoted from Sprite
