package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
)

// Boss3 is the stage 3 boss. Big, aggressive, lots of bullets.
type Boss3 struct {
	*objects.ContainerSprite
	objects.CollisionTester

	HP     int
	MaxHP  int
	isDead bool

	time           float64
	phase          int
	PendingBullets []*EnemyBullet
	flashTimer     int

	// Parts
	Core   *objects.Sprite
	WingL  *objects.Sprite
	WingR  *objects.Sprite
	Cannon *objects.Sprite
}

func NewBoss3(x, y float64) *Boss3 {
	b := &Boss3{}
	b.InitBoss3(x, y)
	return b
}

func (b *Boss3) InitBoss3(x, y float64) {
	// --- Core ---
	// Large Purple Hexagon-ish
	coreImg := ebiten.NewImage(64, 64)
	coreImg.Fill(color.RGBA{100, 0, 150, 255})
	b.ContainerSprite = objects.NewContainerSprite(x, y, coreImg)
	b.Core = b.ContainerSprite.Sprite
	b.PendingBullets = []*EnemyBullet{}

	// --- Wings ---
	wingImg := ebiten.NewImage(32, 96)
	wingImg.Fill(color.RGBA{80, 0, 120, 255})

	b.WingL = objects.NewSprite(-60, 20, wingImg)
	b.AddChild(b.WingL)

	b.WingR = objects.NewSprite(60, 20, wingImg)
	b.AddChild(b.WingR)

	// --- Cannon ---
	cannonImg := ebiten.NewImage(24, 48)
	cannonImg.Fill(color.RGBA{200, 50, 200, 255})
	b.Cannon = objects.NewSprite(0, 40, cannonImg)
	b.AddChild(b.Cannon)

	// --- Attributes ---
	b.MaxHP = 150 // Very Strong
	b.HP = b.MaxHP
	b.phase = 0
	b.flashTimer = 0

	// --- Collision ---
	coreCol := objects.NewRectCollider(b.Core)
	wingLCol := objects.NewRectCollider(b.WingL)
	wingRCol := objects.NewRectCollider(b.WingR)
	cannonCol := objects.NewRectCollider(b.Cannon)

	b.CollisionTester = objects.NewCompositeCollider(
		[]objects.CollisionTester{coreCol, wingLCol, wingRCol, cannonCol},
		objects.CompositeOr,
	)
}

func (b *Boss3) Update() {
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
		b.WingL.ColorScale.Reset()
		b.WingL.ColorScale.Scale(f, f, f, 1)
		b.WingR.ColorScale.Reset()
		b.WingR.ColorScale.Scale(f, f, f, 1)
		b.Cannon.ColorScale.Reset()
		b.Cannon.ColorScale.Scale(f, f, f, 1)
	} else {
		b.ColorScale.Reset()
		b.ColorScale.Scale(1, 1, 1, 1)
		b.WingL.ColorScale.Reset()
		b.WingL.ColorScale.Scale(1, 1, 1, 1)
		b.WingR.ColorScale.Reset()
		b.WingR.ColorScale.Scale(1, 1, 1, 1)
		b.Cannon.ColorScale.Reset()
		b.Cannon.ColorScale.Scale(1, 1, 1, 1)
	}

	// Logic
	if b.phase == 0 {
		if b.Y < 120 {
			b.Y += 1.0
		} else {
			b.phase = 1
		}
	} else {
		// Side to side movement
		b.X = 320 + math.Sin(b.time*0.03)*200

		// Cannon Rotation
		b.Cannon.Angle = math.Sin(b.time*0.1) * 0.3

		// Shooting Pattern
		// 1. Core rapid fire (straight)
		if int(b.time)%80 == 0 {
			bx, by := b.Core.GetGlobalPos()
			b.PendingBullets = append(b.PendingBullets, NewEnemyBullet(bx+16, by+64, math.Pi/2))
			b.PendingBullets = append(b.PendingBullets, NewEnemyBullet(bx+48, by+64, math.Pi/2))
		}

		// 2. Wing spread shot (slow interval)
		if int(b.time)%150 == 0 {
			lx, ly := b.WingL.GetGlobalPos()
			rx, ry := b.WingR.GetGlobalPos()

			// Left wing spread
			for i := -1; i <= 1; i++ {
				angle := math.Pi/2 + float64(i)*0.4
				b.PendingBullets = append(b.PendingBullets, NewEnemyBullet(lx+16, ly+48, angle))
			}
			// Right wing spread
			for i := -1; i <= 1; i++ {
				angle := math.Pi/2 + float64(i)*0.4
				b.PendingBullets = append(b.PendingBullets, NewEnemyBullet(rx+16, ry+48, angle))
			}
		}

		// 3. Cannon aimed shot (if implemented aim, else random spray)
		if int(b.time)%50 == 0 {
			cx, cy := b.Cannon.GetGlobalPos()
			angle := math.Pi/2 + b.Cannon.Angle + (rand.Float64()*0.2 - 0.1)
			b.PendingBullets = append(b.PendingBullets, NewEnemyBullet(cx+12, cy+48, angle))
		}
	}

	b.ContainerSprite.Update()
}

func (b *Boss3) Draw(screen *ebiten.Image) {
	if b.isDead {
		return
	}
	b.ContainerSprite.Draw(screen)
}

func (b *Boss3) IsDead() bool {
	return b.isDead
}

func (b *Boss3) MarkDead() {
	b.isDead = true
}

func (b *Boss3) ApplyDamage(damage int) {
	b.HP -= damage
	b.flashTimer = 4
	if b.HP <= 0 {
		b.MarkDead()
	}
}

func (b *Boss3) DropItem() *Item {
	return nil
}

func (b *Boss3) GetPendingBullets() []*EnemyBullet {
	bullets := b.PendingBullets
	b.PendingBullets = []*EnemyBullet{}
	return bullets
}

func (b *Boss3) GetHP() int {
	return b.HP
}

func (b *Boss3) GetMaxHP() int {
	return b.MaxHP
}
