package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"MyProject/objects"
	"MyProject/parts"
)

// Particle struct as an Entity
type Particle struct {
	parts.EntityBase
	VX, VY  float64
	Color   color.Color
	Radius  float32
	Life    float64
	MaxLife float64
	Alpha   float32
}

func NewParticle(x, y float64) *Particle {
	p := &Particle{}
	p.InitParticle(x, y)
	return p
}

func (p *Particle) InitParticle(x, y float64) {
	p.EntityBase.InitEntityBase(p, x, y)
}

func (p *Particle) Update() {
	p.X += p.VX
	p.Y += p.VY
	p.Life--
	p.Alpha = float32(p.Life) / float32(p.MaxLife)

	if p.Life <= 0 {
		// Remove from parent
		if parent := p.EntityBase.Parent; parent != nil {
			// Need to cast parent to something that can remove child.
			// EffectManager embeds Container, which implements Grouping.
			// Grouping has RemoveChild.
			// But Parent is stored as Entity.
			// We can check if Parent implements a interface with RemoveChild,
			// or assume it's a *EffectManager or *objects.Container.

			// Checking parts.Grouping would be ideal if exported, or checking if it has RemoveChild.
			// Since Grouping struct is in parts package, and RemoveChild takes Entity.

			// Use type assertion for now. EffectManager embeds objects.Container.
			if container, ok := parent.(*EffectManager); ok {
				container.RemoveChild(p)
			} else if container, ok := parent.(*objects.Container); ok {
				container.RemoveChild(p)
			}
		}
	}
}

func (p *Particle) Draw(screen *ebiten.Image) {
	// Calculate color with alpha
	r, g, b, _ := p.Color.RGBA()

	a := uint8(p.Alpha * 255)
	c := color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: a,
	}

	// GetGlobalPos() is available via EntityBase, but vector.DrawFilledCircle takes screen coords?
	// The screen passed to Draw is usually the target image.
	// If we are children of a Container that handles transforms (Scale/Rotate),
	// Ebiten's standard DrawImage handles matrix.
	// BUT vector.DrawFilledCircle draws directly to the image at x, y.
	// It does NOT respect Ebiten's DrawImageOptions (GeoM) because it's a direct pixel operation helper.

	// Issue: If the parent Container moves/rotates, vector.DrawFilledCircle won't automatically transform
	// unless we calculate global position manually.
	// Container's Draw implementation calls child.Draw(screen). It doesn't apply transform to screen.
	// Wait, objects.Container uses Grouping. Grouping draws children.
	// Usually in Ebiten, hierarchy is handled by passing GeoM to DrawImage.
	// But vector package draws primitives.

	// If the EffectManager is at 0,0 and never moves, then local Pos = Global Pos.
	// If EffectManager moves, we need Global Pos.

	gx, gy := p.GetGlobalPos()
	vector.DrawFilledCircle(screen, float32(gx), float32(gy), p.Radius, c, true)
}

// EffectManager embeds Container
type EffectManager struct {
	*objects.Container
}

func NewEffectManager() *EffectManager {
	em := &EffectManager{}
	em.InitEffectManager()
	return em
}

func (em *EffectManager) InitEffectManager() {
	em.Container = objects.NewContainer(0, 0)
	// We init the embedded container, which sets up EntityBase and Grouping
}

// Update and Draw are handled by Container (Grouping)

func (em *EffectManager) SpawnExplosion(centerx, centery float64, baseColor color.Color) {
	count := 10 + rand.Intn(10) // 10 to 20 particles
	for i := 0; i < count; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := 1.0 + rand.Float64()*3.0

		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		life := 20.0 + rand.Float64()*20.0 // 20-40 frames

		p := NewParticle(centerx, centery)
		p.VX = vx
		p.VY = vy
		p.Color = baseColor
		p.Radius = 2.0 + rand.Float32()*2.0
		p.Life = life
		p.MaxLife = life
		p.Alpha = 1.0

		em.AddChild(p)
	}
}

func (em *EffectManager) SpawnBossExplosion(centerx, centery float64) {
	// Flashy explosion: Multiple rings/bursts
	colors := []color.Color{
		color.RGBA{255, 50, 50, 255},   // Red
		color.RGBA{255, 150, 50, 255},  // Orange
		color.RGBA{255, 255, 50, 255},  // Yellow
		color.RGBA{200, 200, 200, 255}, // White/Smoke
	}

	for i := 0; i < 100; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := 2.0 + rand.Float64()*6.0

		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		life := 40.0 + rand.Float64()*60.0 // 40-100 frames

		p := NewParticle(centerx, centery)
		p.VX = vx
		p.VY = vy
		p.Color = colors[rand.Intn(len(colors))]
		p.Radius = 3.0 + rand.Float32()*4.0
		p.Life = life
		p.MaxLife = life
		p.Alpha = 1.0

		em.AddChild(p)
	}
}
