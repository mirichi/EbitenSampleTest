package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"MyProject/objects"
)

// Background scrolls vertically
type Background struct {
	*objects.Sprite
	ScrollRate float64
	offsetY    float64
	ParallaxX  float64
}

func NewBackground(imgPath string, scrollRate float64) *Background {
	img, _, err := ebitenutil.NewImageFromFile(imgPath)
	if err != nil {
		// Fallback specific for this game concept if file load fails.
		// In production, we should handle this better.
		// For now, create a fallback image.
		img = ebiten.NewImage(ScreenWidth, ScreenHeight)
		// img.Fill(color.Black) // Default black
	}

	// We need 2 images to loop seamless?
	// Or just draw the same image twice?
	// Ebiten's infinite scrolling usually involves drawing the same image twice.
	// But `Sprite` draws one image.
	// We can implement Custom Draw in Background struct.

	bg := &Background{
		Sprite:     objects.NewSprite(0, 0, img),
		ScrollRate: scrollRate,
		offsetY:    0,
		ParallaxX:  0,
	}
	// Reset Origin to Top-Left for easier manual tiling
	bg.OriginX = 0
	bg.OriginY = 0
	// Center horizontally
	w := float64(img.Bounds().Dx())
	bg.X = (ScreenWidth - w) / 2

	return bg
}

func (bg *Background) Update() {
	bg.offsetY += bg.ScrollRate
	h := float64(bg.Sprite.Image.Bounds().Dy())
	if bg.offsetY >= h {
		bg.offsetY -= h
	}
}

func (bg *Background) Draw(screen *ebiten.Image) {
	w := float64(bg.Sprite.Image.Bounds().Dx())
	h := float64(bg.Sprite.Image.Bounds().Dy())

	// Calculate base X with wrapping
	// We want to cover the screen horizontally.
	// ParallaxX is the offset.
	x := bg.ParallaxX
	// Normalize x to range (-w, 0] to ensure consistent tiling start
	for x > 0 {
		x -= w
	}
	for x <= -w {
		x += w
	}

	// Tiling positions
	// Vertical: Current (offsetY) and Above (offsetY - h)
	ys := []float64{bg.offsetY, bg.offsetY - h}
	// Horizontal: Left (x-w), Current (x) and Right (x + w)
	// Extra tile to left ensures coverage when shifting right
	xs := []float64{x - w, x, x + w}

	op := &ebiten.DrawImageOptions{}

	// Global Matrix (Parent Transform)
	gm := bg.GetGlobalMatrix()

	for _, drawY := range ys {
		for _, drawX := range xs {
			op.GeoM.Reset()
			op.GeoM.Translate(drawX, drawY)
			op.GeoM.Concat(gm) // Apply World Transform
			screen.DrawImage(bg.Sprite.Image, op)
		}
	}
}

// Override Sprite's methods if necessary, or just not add it as simple child if we custom draw.
// But Scene uses `objects.Container`. Containers verify `Draw`.
// `Background` implements `Draw(screen)` so it's `Drawable`.
// It implements `Update` so it's `Updatable` (Entity).
// We should ensure it satisfies `parts.Entity`.
// `*objects.Sprite` satisfies `Entity`.
// `Background` embeds `*Sprite`, so it satisfies `Entity`.
// `Draw` is overridden. Good.
