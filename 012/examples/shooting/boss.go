package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"MyProject/examples/shooting/systems"
	"MyProject/objects"
)

type Boss struct {
	*objects.ContainerSprite
	// Body removed, Boss itself is the body
	LeftArm   *objects.Container
	RightArm  *objects.Container
	LeftHand  *objects.Sprite
	RightHand *objects.Sprite

	// Refs
	EffectManager *EffectManager

	objects.CollisionTester

	// HitBoxes
	BodyCol      objects.CollisionTester
	LeftHandCol  objects.CollisionTester
	RightHandCol objects.CollisionTester

	HP     int
	MaxHP  int
	isDead bool

	time           float64
	phaseTime      float64 // Added for smooth phase transition
	phase          int     // 0: Entering, 1: Fighting
	PendingBullets []*EnemyBullet
	flashTimer     int

	// Hand Attributes
	LeftHandHP    int
	RightHandHP   int
	LeftHandDead  bool
	RightHandDead bool
}

func NewBoss(x, y float64, em *EffectManager) *Boss {
	b := &Boss{}
	b.InitBoss(x, y, em)
	return b
}

func (b *Boss) InitBoss(x, y float64, em *EffectManager) {
	b.EffectManager = em
	// --- Body (Boss itself) ---
	bodyImg := systems.GetResourceManager().GetImage(systems.ImgBossBody)

	b.ContainerSprite = objects.NewContainerSprite(x, y, bodyImg)
	b.PendingBullets = []*EnemyBullet{}

	// --- Arms (Containers for rotation) ---
	// Left Arm
	b.LeftArm = objects.NewContainer(-40, 0)
	// Left Hand Sprite
	handImg := systems.GetResourceManager().GetImage(systems.ImgBossHand)
	b.LeftHand = objects.NewSprite(18, 62, handImg) // Pivot at wrist/connection
	b.LeftHand.OriginX = 18
	b.LeftHand.OriginY = 0
	b.LeftArm.AddChild(b.LeftHand)
	b.AddChild(b.LeftArm)

	// Right Arm
	b.RightArm = objects.NewContainer(40, 0)
	b.RightHand = objects.NewSprite(18, 62, handImg)
	b.RightHand.OriginX = 18
	b.RightHand.OriginY = 0
	b.RightArm.AddChild(b.RightHand)
	b.AddChild(b.RightArm)

	// --- Attributes ---
	b.MaxHP = 70
	b.HP = b.MaxHP
	b.phase = 0 // Entering
	b.LeftHandHP = 20
	b.RightHandHP = 20

	// --- Collision ---
	// --- Collision ---
	// Composite collider: Body + Hands
	// Use b.ContainerSprite.Sprite for the collision tester
	b.BodyCol = objects.NewRectCollider(&b.Sprite)
	b.LeftHandCol = objects.NewCircleCollider(b.LeftHand, objects.Vector2{X: 18, Y: 18}, 18)
	b.RightHandCol = objects.NewCircleCollider(b.RightHand, objects.Vector2{X: 18, Y: 18}, 18)

	b.CollisionTester = objects.NewCompositeCollider(
		[]objects.CollisionTester{b.BodyCol, b.LeftHandCol, b.RightHandCol},
		objects.CompositeOr,
	)
}

func (b *Boss) Update() {
	if b.isDead {
		return
	}
	b.time += 1.0

	// Left Hand dead visibility
	if b.LeftHandDead {
		b.LeftHand.ColorScale.Reset()
		b.LeftHand.ColorScale.Scale(0.5, 0.5, 0.5, 0.5) // Dark and transparent if dead
	}

	// Right Hand dead visibility
	if b.RightHandDead {
		b.RightHand.ColorScale.Reset()
		b.RightHand.ColorScale.Scale(0.5, 0.5, 0.5, 0.5)
	}

	// Logic
	if b.phase == 0 {
		// Moving down to position
		if b.Y < 100 {
			b.Y += 2
		} else {
			b.phase = 1
			b.phaseTime = 0
		}
	} else {
		b.phaseTime += 1.0
		// Bobbing
		b.Y = 100 + math.Sin(b.phaseTime*0.05)*10
		b.X = 320 + math.Sin(b.phaseTime*0.02)*100

		// Shooting logic (every 60 frames approx)
		if int(b.phaseTime)%60 == 0 {
			// Left Hand Shot
			if !b.LeftHandDead {
				lx, ly := b.LeftHand.GetGlobalPos()
				lb := NewEnemyBullet(lx, ly+18, math.Pi/2+rand.Float64()*0.5-0.25)
				b.PendingBullets = append(b.PendingBullets, lb)
			}

			// Right Hand Shot
			if !b.RightHandDead {
				rx, ry := b.RightHand.GetGlobalPos()
				rb := NewEnemyBullet(rx, ry+18, math.Pi/2+rand.Float64()*0.5-0.25)
				b.PendingBullets = append(b.PendingBullets, rb)
			}
		}
	}

	// Arm Rotation
	b.LeftArm.Angle = math.Sin(b.time*0.05) * 0.5
	b.RightArm.Angle = -math.Sin(b.time*0.05) * 0.5

	// Hand waving
	b.LeftHand.Angle = math.Cos(b.time*0.1) * 0.5
	b.RightHand.Angle = -math.Cos(b.time*0.1) * 0.5

	b.ContainerSprite.Update()
}

func (b *Boss) Draw(screen *ebiten.Image) {
	if b.isDead {
		return
	}
	b.ContainerSprite.Draw(screen)

	// Draw Hand HP Bars
	if !b.LeftHandDead {
		lx, ly := b.LeftHand.GetGlobalPos()
		// Bar Size: 30x4
		// Position: Below hand (approx 20px down from center)
		// MaxHP: 20
		ratio := float32(b.LeftHandHP) / 20.0
		// Background (Gray)
		vector.DrawFilledRect(screen, float32(lx-15), float32(ly+30), 30, 4, color.RGBA{50, 50, 50, 255}, true)
		// HP (Green)
		vector.DrawFilledRect(screen, float32(lx-15), float32(ly+30), 30*ratio, 4, color.RGBA{0, 255, 0, 255}, true)
	}

	if !b.RightHandDead {
		rx, ry := b.RightHand.GetGlobalPos()
		ratio := float32(b.RightHandHP) / 20.0
		vector.DrawFilledRect(screen, float32(rx-15), float32(ry+30), 30, 4, color.RGBA{50, 50, 50, 255}, true)
		vector.DrawFilledRect(screen, float32(rx-15), float32(ry+30), 30*ratio, 4, color.RGBA{0, 255, 0, 255}, true)
	}
}

func (b *Boss) IsDead() bool {
	return b.isDead
}

func (b *Boss) MarkDead() {
	b.isDead = true
}

func (b *Boss) ApplyDamage(damage int) {
	b.HP -= damage
	b.HP -= damage
	b.Sprite.Flash(4)
	if b.HP <= 0 {
		b.MarkDead()
	}
}

func (b *Boss) DropItem() *Item {
	return nil
}

func (b *Boss) GetPendingBullets() []*EnemyBullet {
	bullets := b.PendingBullets
	b.PendingBullets = []*EnemyBullet{}
	return bullets
}

func (b *Boss) GetHP() int {
	return b.HP
}

func (b *Boss) GetMaxHP() int {
	return b.MaxHP
}

func (b *Boss) HandleHit(bullet *Bullet) {
	// Use specific colliders for accurate hit detection
	if !b.LeftHandDead && bullet.Test(b.LeftHandCol) {
		b.LeftHandHP--
		b.LeftHand.Flash(4)
		if b.LeftHandHP <= 0 {
			b.LeftHandDead = true
			if b.EffectManager != nil {
				bx, by := b.GetGlobalPos()
				lhx, lhy := b.LeftHand.GetGlobalPos()
				diffX := lhx - bx
				diffY := lhy - by

				spawnX := b.X + diffX
				spawnY := b.Y + diffY

				b.EffectManager.SpawnExplosion(spawnX, spawnY, color.RGBA{150, 150, 0, 255})
			}
		}
		return
	}

	if !b.RightHandDead && bullet.Test(b.RightHandCol) {
		b.RightHandHP--
		b.RightHand.Flash(4)
		if b.RightHandHP <= 0 {
			b.RightHandDead = true
			if b.EffectManager != nil {
				bx, by := b.GetGlobalPos()
				hx, hy := b.RightHand.GetGlobalPos()
				spawnX := b.X + (hx - bx)
				spawnY := b.Y + (hy - by)
				b.EffectManager.SpawnExplosion(spawnX, spawnY, color.RGBA{150, 150, 0, 255})
			}
		}
		return
	}

	if bullet.Test(b.BodyCol) {
		b.ApplyDamage(1)
	}
}
