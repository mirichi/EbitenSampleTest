package main

import (
	"image/color"
	"math"

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

	time  float64
	phase int // 0: Entering, 1: Fighting
}

func NewBoss(x, y float64) *Boss {
	b := &Boss{}
	b.InitBoss(x, y)
	return b
}

func (b *Boss) InitBoss(x, y float64) {
	b.EntityBase.InitEntityBase(b, x, y)
	b.Root = objects.NewContainer(x, y)

	// --- Body ---
	bodyImg := ebiten.NewImage(64, 64)
	bodyImg.Fill(color.RGBA{200, 50, 50, 255}) // Reddish body
	b.Body = objects.NewSprite(0, 0, bodyImg)
	b.Body.OriginX = 32
	b.Body.OriginY = 32
	b.Root.AddChild(b.Body)

	// --- Arms (Containers for rotation) ---
	// Left Arm
	b.LeftArm = objects.NewContainer(-40, 0)
	// Left Hand Sprite
	handImg := ebiten.NewImage(24, 24)
	handImg.Fill(color.RGBA{150, 150, 0, 255})     // Yellowish hand
	b.LeftHand = objects.NewSprite(0, 30, handImg) // Offset from arm joint
	b.LeftHand.OriginX = 12
	b.LeftHand.OriginY = 0 // Pivot at wrist/connection
	b.LeftArm.AddChild(b.LeftHand)
	b.Root.AddChild(b.LeftArm)

	// Right Arm
	b.RightArm = objects.NewContainer(40+34, 0)
	b.RightHand = objects.NewSprite(0, 30, handImg)
	b.RightHand.OriginX = 12
	b.RightHand.OriginY = 0
	b.RightArm.AddChild(b.RightHand)
	b.Root.AddChild(b.RightArm)

	// --- Attributes ---
	b.MaxHP = 50
	b.HP = b.MaxHP
	b.phase = 0 // Entering

	// --- Collision ---
	// Composite collider: Body + Hands
	// Note: We use the Sprites directly for primitive colliders
	bodyCol := objects.NewRectCollider(b.Body)
	lHandCol := objects.NewCircleCollider(b.LeftHand, objects.Vector2{X: 12, Y: 12}, 12)
	rHandCol := objects.NewCircleCollider(b.RightHand, objects.Vector2{X: 12, Y: 12}, 12)

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
