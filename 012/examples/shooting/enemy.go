package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
	"MyProject/parts"
)

// Enemy Interface
type Enemy interface {
	parts.Entity
	Update()
	Draw(screen *ebiten.Image)
	objects.CollisionTester
	IsDead() bool
	MarkDead()
	ApplyDamage(damage int)
	DropItem() *Item
	GetPendingBullets() []*EnemyBullet
	GetHP() int
	GetMaxHP() int
}

// EnemyBase Struct
type EnemyBase struct {
	*objects.Sprite // Embed Sprite
	*objects.PolygonCollider
	isDead bool

	// Common state for logic
	baseX  float64
	tick   int
	speedY float64
	hp     int
	maxHp  int
}

// ...
func (e *EnemyBase) GetHP() int {
	return e.hp
}

func (e *EnemyBase) GetMaxHP() int {
	return e.maxHp // Default 0/1 depending on init
}

func (e *EnemyBase) Draw(screen *ebiten.Image) {
	if !e.isDead {
		e.Sprite.Draw(screen)
	}
}

func (e *EnemyBase) IsDead() bool {
	return e.isDead
}

func (e *EnemyBase) MarkDead() {
	e.isDead = true
}

func (e *EnemyBase) ApplyDamage(damage int) {
	e.hp -= damage
	if e.hp <= 0 {
		e.MarkDead()
	}
}

func (e *EnemyBase) DropItem() *Item {
	return nil
}

func (e *EnemyBase) GetPendingBullets() []*EnemyBullet {
	return nil
}

func (e *EnemyBase) checkBounds() {
	if e.Sprite.Y > ScreenHeight+32 {
		e.isDead = true
	}
}

// --- Concrete Implementations ---

// StraightEnemy
type StraightEnemy struct {
	EnemyBase
}

func (e *StraightEnemy) Update() {
	if e.isDead {
		return
	}
	e.tick++

	e.Sprite.Y += e.speedY
	e.Sprite.Angle += 0.05

	e.Sprite.Update()
	e.checkBounds()
}

// WaveEnemy
type WaveEnemy struct {
	EnemyBase
}

func (e *WaveEnemy) Update() {
	if e.isDead {
		return
	}
	e.tick++

	e.Sprite.Y += e.speedY
	// Sine wave movement
	e.Sprite.X = e.baseX + math.Sin(float64(e.tick)*0.05)*50
	e.Sprite.Angle += 0.02

	e.Sprite.Update()
	e.checkBounds()
}

// FastEnemy
type FastEnemy struct {
	EnemyBase
}

func (e *FastEnemy) Update() {
	if e.isDead {
		return
	}
	e.tick++

	e.Sprite.Y += e.speedY
	// Zigzag
	if (e.tick/30)%2 == 0 {
		e.Sprite.X += 1
	} else {
		e.Sprite.X -= 1
	}
	e.Sprite.Angle += 0.1

	e.Sprite.Update()
	e.checkBounds()
}

// MediumEnemy (Slow, Hard, Drops Item)
type MediumEnemy struct {
	EnemyBase
}

func (e *MediumEnemy) Update() {
	if e.isDead {
		return
	}
	e.tick++

	e.Sprite.Y += e.speedY
	e.Sprite.Angle -= 0.02 // Reverse rotation

	e.Sprite.Update()
	e.checkBounds()
}

func (e *MediumEnemy) DropItem() *Item {
	return NewItem(e.Sprite.X, e.Sprite.Y, ItemTypePowerUp)
}

// ShootingEnemy (Moves and Shoots)
type ShootingEnemy struct {
	EnemyBase
	PendingBullets []*EnemyBullet
}

func (e *ShootingEnemy) Update() {
	if e.isDead {
		return
	}
	e.tick++

	e.Sprite.Y += e.speedY
	e.Sprite.Angle += 0.01

	// Shooting (every 60 frames)
	if e.tick%120 == 0 {
		// 3-way shot downwards
		bx, by := e.Sprite.X, e.Sprite.Y
		b1 := NewEnemyBullet(bx, by, math.Pi/2)
		b2 := NewEnemyBullet(bx, by, math.Pi/2-0.3)
		b3 := NewEnemyBullet(bx, by, math.Pi/2+0.3)
		e.PendingBullets = append(e.PendingBullets, b1, b2, b3)
	}

	e.Sprite.Update()
	e.checkBounds()
}

func (e *ShootingEnemy) GetPendingBullets() []*EnemyBullet {
	bullets := e.PendingBullets
	e.PendingBullets = []*EnemyBullet{}
	return bullets
}
