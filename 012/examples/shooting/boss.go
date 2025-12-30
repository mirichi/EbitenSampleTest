package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
	"MyProject/parts"
)

type Boss struct {
	parts.EntityBase
	Root      *objects.Container
	Body      *objects.Sprite
	LeftArm   *objects.Container
	RightArm  *objects.Container
	LeftHand  *objects.Sprite
	RightHand *objects.Sprite

	Collider objects.CollisionTester

	HP     int
	MaxHP  int
	isDead bool

	time           float64
	phase          int // 0: Entering, 1: Fighting
	PendingBullets []*EnemyBullet
	flashTimer     int
}

func NewBoss(x, y float64) *Boss {
	b := &Boss{}
	b.InitBoss(x, y)
	return b
}

func (b *Boss) InitBoss(x, y float64) {
	b.EntityBase.InitEntityBase(b, x, y)
	b.Root = objects.NewContainer(x, y)
	b.PendingBullets = []*EnemyBullet{}

	// --- Body ---
	bodyImg := ebiten.NewImage(96, 96)
	bodyImg.Fill(color.RGBA{200, 50, 50, 255}) // Reddish body
	b.Body = objects.NewSprite(0, 0, bodyImg)
	b.Root.AddChild(b.Body)

	// --- Arms (Containers for rotation) ---
	// Left Arm
	b.LeftArm = objects.NewContainer(-40, 0)
	// Left Hand Sprite
	handImg := ebiten.NewImage(36, 36)
	handImg.Fill(color.RGBA{150, 150, 0, 255})      // Yellowish hand
	b.LeftHand = objects.NewSprite(18, 62, handImg) // Pivot at wrist/connection
	b.LeftHand.OriginX = 18
	b.LeftHand.OriginY = 0
	b.LeftArm.AddChild(b.LeftHand)
	b.Root.AddChild(b.LeftArm)

	// Right Arm
	b.RightArm = objects.NewContainer(40, 0)
	b.RightHand = objects.NewSprite(18, 62, handImg)
	b.RightHand.OriginX = 18
	b.RightHand.OriginY = 0
	b.RightArm.AddChild(b.RightHand)
	b.Root.AddChild(b.RightArm)

	// --- Attributes ---
	b.MaxHP = 50
	b.HP = b.MaxHP
	b.phase = 0 // Entering
	b.flashTimer = 0

	// --- Collision ---
	// Composite collider: Body + Hands
	// Note: We use the Sprites directly for primitive colliders
	bodyCol := objects.NewRectCollider(b.Body)
	lHandCol := objects.NewCircleCollider(b.LeftHand, objects.Vector2{X: 18, Y: 18}, 18)
	rHandCol := objects.NewCircleCollider(b.RightHand, objects.Vector2{X: 18, Y: 18}, 18)

	b.Collider = objects.NewCompositeCollider(
		[]objects.CollisionTester{bodyCol, lHandCol, rHandCol},
		objects.CompositeOr,
	)
}

func (b *Boss) Update() {
	if b.isDead {
		return
	}
	b.time += 1.0

	// Flash Logic
	if b.flashTimer > 0 {
		b.flashTimer--
		// Flash Effect: Scale color to be very bright (approximating white flash)
		f := float32(2.0)
		b.Body.ColorScale.Reset()
		b.Body.ColorScale.Scale(f, f, f, 1)
		b.LeftHand.ColorScale.Reset()
		b.LeftHand.ColorScale.Scale(f, f, f, 1)
		b.RightHand.ColorScale.Reset()
		b.RightHand.ColorScale.Scale(f, f, f, 1)
	} else {
		// Normal Color
		b.Body.ColorScale.Reset()
		b.Body.ColorScale.Scale(1, 1, 1, 1)
		b.LeftHand.ColorScale.Reset()
		b.LeftHand.ColorScale.Scale(1, 1, 1, 1)
		b.RightHand.ColorScale.Reset()
		b.RightHand.ColorScale.Scale(1, 1, 1, 1)
	}

	// Logic
	if b.phase == 0 {
		// Moving down to position
		if b.Root.Y < 100 {
			b.Root.Y += 2
		} else {
			b.phase = 1
		}
	} else {
		// Bobbing
		b.Root.Y = 100 + math.Sin(b.time*0.05)*10
		b.Root.X = 320 + math.Cos(b.time*0.02)*100

		// Shooting logic (every 60 frames approx)
		if int(b.time)%60 == 0 {
			// Left Hand Shot
			lx, ly := b.LeftHand.GetGlobalPos()
			lb := NewEnemyBullet(lx, ly+18, math.Pi/2+rand.Float64()*0.5-0.25)
			b.PendingBullets = append(b.PendingBullets, lb)

			// Right Hand Shot
			rx, ry := b.RightHand.GetGlobalPos()
			rb := NewEnemyBullet(rx, ry+18, math.Pi/2+rand.Float64()*0.5-0.25)
			b.PendingBullets = append(b.PendingBullets, rb)
		}
	}

	// Arm Rotation
	b.LeftArm.Angle = math.Sin(b.time*0.05) * 0.5
	b.RightArm.Angle = -math.Sin(b.time*0.05) * 0.5

	// Hand waving
	b.LeftHand.Angle = math.Cos(b.time*0.1) * 0.5
	b.RightHand.Angle = -math.Cos(b.time*0.1) * 0.5

	b.Root.Update()
}

func (b *Boss) Draw(screen *ebiten.Image) {
	if b.isDead {
		return
	}
	b.Root.Draw(screen)
}

func (b *Boss) IsDead() bool {
	return b.isDead
}

func (b *Boss) MarkDead() {
	b.HP--
	b.flashTimer = 4 // Flash for 4 frames
	if b.HP <= 0 {
		b.isDead = true
	}
}

// Implement CollisionTester interface for Boss
func (b *Boss) Test(c objects.CollisionTester) bool {
	return b.Collider.Test(c)
}

func (b *Boss) CollisionShape() objects.CollisionTester {
	return b.Collider
}
