package gameparts

import (
	"MyProject/parts"

	"github.com/hajimehoshi/ebiten/v2"
)

// Flasher manages the flash effect for sprites.
type Flasher struct {
	e        parts.Entity
	timer    int
	maxTimer int
	target   *ebiten.ColorScale
}

// NewFlasher creates a new Flasher for a target ColorScale.
func NewFlasher(e parts.Entity, target *ebiten.ColorScale) *Flasher {
	f := &Flasher{}
	f.InitFlasher(e, target)
	return f
}

func (f *Flasher) InitFlasher(e parts.Entity, target *ebiten.ColorScale) {
	f.e = e
	f.target = target
	f.e.GetEntityBase().AddAfterUpdateFunction(f.UpdateFunction)
}

// Flash triggers the flash effect for the specified duration (frames).
func (f *Flasher) Flash(frames int) {
	f.timer = frames
	f.maxTimer = frames
}

// UpdateFunction updates the flash state. Should be called every frame.
func (f *Flasher) UpdateFunction() {
	if f.timer > 0 {
		f.timer--
		// Flash Effect: simple white flash logic
		// You can customize the curve or intensity here
		scale := float32(2.0)
		f.target.Reset()
		f.target.Scale(scale, scale, scale, 1)
	} else if f.timer == 0 {
		// Just finished, reset
		f.target.Reset()
		f.timer = -1 // Mark as idle
	}
}

// IsFlashing returns true if currently flashing.
func (f *Flasher) IsFlashing() bool {
	return f.timer >= 0
}
