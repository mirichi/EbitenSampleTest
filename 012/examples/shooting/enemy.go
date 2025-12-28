package main

import (
	"image/color"
	"math"
	"math/rand"

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
}

// EnemyBase Struct
type EnemyBase struct {
	*objects.Sprite // Embed Sprite
	*objects.PolygonCollider
	isDead bool

	// Common state for logic
	baseX  float64
	time   float64
	speedY float64
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
	e.time += 1.0

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
	e.time += 1.0

	e.Sprite.Y += e.speedY
	// Sine wave movement
	e.Sprite.X = e.baseX + math.Sin(e.time*0.05)*50
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
	e.time += 1.0

	e.Sprite.Y += e.speedY
	// Zigzag
	if int(e.time/30)%2 == 0 {
		e.Sprite.X += 1
	} else {
		e.Sprite.X -= 1
	}
	e.Sprite.Angle += 0.1

	e.Sprite.Update()
	e.checkBounds()
}

// Factory
func SpawnEnemy() Enemy {
	t := rand.Intn(3)

	var img *ebiten.Image
	var speed float64
	var col color.Color
	w, h := 32, 32

	switch t {
	case 0: // Straight
		col = color.RGBA{255, 0, 0, 255}
		speed = 2.0
	case 1: // Wave
		col = color.RGBA{0, 255, 0, 255}
		speed = 1.5
	case 2: // Fast
		col = color.RGBA{255, 255, 0, 255}
		speed = 4.0
	}

	img = ebiten.NewImage(w, h)
	img.Fill(col)

	startX := rand.Float64()*(ScreenWidth-float64(w)) + float64(w)/2

	s := objects.NewSprite(startX, -float64(h), img)
	s.OriginX = float64(w) / 2
	s.OriginY = float64(h) / 2
	s.Angle = rand.Float64() * 6.28

	base := EnemyBase{
		Sprite: s,
		isDead: false,
		baseX:  startX,
		speedY: speed,
	}
	base.PolygonCollider = objects.NewRectCollider(s)

	switch t {
	case 0:
		return &StraightEnemy{EnemyBase: base}
	case 1:
		return &WaveEnemy{EnemyBase: base}
	case 2:
		return &FastEnemy{EnemyBase: base}
	default:
		return &StraightEnemy{EnemyBase: base}
	}
}
