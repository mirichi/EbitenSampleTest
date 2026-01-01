package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
)

// Boss2 is a new boss with different behavior and appearance
type Boss2 struct {
	*Boss // Inherit basic structure if possible, but composition is cleaner.
	// Since Boss struct has specific fields like Arms, maybe we should make a BaseBoss or just copy logic.
	// Copy logic for now to allow distinct behavior.

	*objects.ContainerSprite
	objects.CollisionTester

	HP     int
	MaxHP  int
	isDead bool

	time           float64
	phaseTime      float64
	phase          int
	PendingBullets []*EnemyBullet
	flashTimer     int

	// Specific parts
	MainBody *objects.Sprite
	Turret   *objects.Sprite // Rotating turret
}

func NewBoss2(x, y float64) *Boss2 {
	b := &Boss2{}
	b.InitBoss2(x, y)
	return b
}

func (b *Boss2) InitBoss2(x, y float64) {
	// --- Body ---
	// Blueish tank-like body
	bodyImg := ebiten.NewImage(120, 80)
	bodyImg.Fill(color.RGBA{50, 50, 200, 255})

	b.ContainerSprite = objects.NewContainerSprite(x, y, bodyImg)
	b.MainBody = b.ContainerSprite.Sprite
	b.PendingBullets = []*EnemyBullet{}

	// --- Turret ---
	turretImg := ebiten.NewImage(40, 40)
	turretImg.Fill(color.RGBA{100, 200, 255, 255})
	b.Turret = objects.NewSprite(0, 0, turretImg)
	b.Turret.OriginX = 20
	b.Turret.OriginY = 20
	b.AddChild(b.Turret)

	// --- Attributes ---
	b.MaxHP = 80 // Stronger
	b.HP = b.MaxHP
	b.phase = 0
	b.flashTimer = 0

	// --- Collision ---
	bodyCol := objects.NewRectCollider(b.MainBody)
	turretCol := objects.NewCircleCollider(b.Turret, objects.Vector2{X: 20, Y: 20}, 20)

	b.CollisionTester = objects.NewCompositeCollider(
		[]objects.CollisionTester{bodyCol, turretCol},
		objects.CompositeOr,
	)
}

func (b *Boss2) Update() {
	if b.isDead {
		return
	}
	b.time += 1.0

	// Flash Logic
	if b.flashTimer > 0 {
		b.flashTimer--
		f := float32(2.0)
		b.ColorScale.Reset()
		b.ColorScale.Scale(f, f, f, 1)
		b.Turret.ColorScale.Reset()
		b.Turret.ColorScale.Scale(f, f, f, 1)
	} else {
		b.ColorScale.Reset()
		b.ColorScale.Scale(1, 1, 1, 1)
		b.Turret.ColorScale.Reset()
		b.Turret.ColorScale.Scale(1, 1, 1, 1)
	}

	// Movement
	if b.phase == 0 {
		if b.Y < 80 {
			b.Y += 1.5
		} else {
			b.phase = 1
			b.phaseTime = 0
		}
	} else {
		b.phaseTime += 1.0
		// Figure-8 movement
		b.X = 320 + math.Sin(b.phaseTime*0.02)*150
		b.Y = 80 + math.Sin(b.phaseTime*0.04)*50

		// Turret Rotation
		b.Turret.Angle += 0.05

		// Shooting (Spread shot)
		if int(b.phaseTime)%90 == 0 {
			tx, ty := b.Turret.GetGlobalPos()

			// 5-way shot
			for i := -2; i <= 2; i++ {
				angle := math.Pi/2 + float64(i)*0.2 + b.Turret.Angle // Add turret rotation to shot
				bullet := NewEnemyBullet(tx, ty, angle)
				b.PendingBullets = append(b.PendingBullets, bullet)
			}
		}
	}

	b.ContainerSprite.Update()
}

func (b *Boss2) Draw(screen *ebiten.Image) {
	if b.isDead {
		return
	}
	b.ContainerSprite.Draw(screen)
}

func (b *Boss2) IsDead() bool {
	return b.isDead
}

func (b *Boss2) MarkDead() {
	b.isDead = true
}

func (b *Boss2) ApplyDamage(damage int) {
	b.HP -= damage
	b.flashTimer = 4
	if b.HP <= 0 {
		b.MarkDead()
	}
}

func (b *Boss2) DropItem() *Item {
	return nil
}

func (b *Boss2) GetPendingBullets() []*EnemyBullet {
	bullets := b.PendingBullets
	b.PendingBullets = []*EnemyBullet{}
	return bullets
}

func (b *Boss2) GetHP() int {
	return b.HP
}

func (b *Boss2) GetMaxHP() int {
	return b.MaxHP
}

// Ensure Boss2 implements Enemy interface (which requires parts.Entity etc.)
// Boss2 embeds ContainerSprite -> Sprite -> EntityBase, so yes.
