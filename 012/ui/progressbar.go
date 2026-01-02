package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ProgressBar is a UI component that displays progress or a value within a range.
type ProgressBar struct {
	BlankControl
	Min, Max, Value float64
	FillColor       color.Color
	BackgroundColor color.Color
	BorderColor     color.Color
}

// NewProgressBar creates a new ProgressBar.
func NewProgressBar(x, y, w, h float64) *ProgressBar {
	p := &ProgressBar{}
	p.InitProgressBar(x, y, w, h)
	return p
}

func (p *ProgressBar) InitProgressBar(x, y, w, h float64) {
	p.InitBlankControl(x, y, w, h)
	p.Min = 0
	p.Max = 100
	p.Value = 50
	p.FillColor = color.RGBA{0, 200, 0, 255}        // Green
	p.BackgroundColor = color.RGBA{50, 50, 50, 255} // Dark Gray
	p.BorderColor = color.White

	p.AddBeforeDrawFunction(p.DrawFunction)
}

// SetRange sets the minimum and maximum values of the progress bar.
func (p *ProgressBar) SetRange(min, max float64) {
	p.Min = min
	p.Max = max
}

// SetValue sets the current value of the progress bar.
func (p *ProgressBar) SetValue(v float64) {
	if v < p.Min {
		v = p.Min
	}
	if v > p.Max {
		v = p.Max
	}
	p.Value = v
}

// Draw renders the progress bar.
func (p *ProgressBar) DrawFunction(screen *ebiten.Image) {
	if !p.Visible {
		return
	}

	x, y := p.GetGlobalPos()
	w, h := p.Width, p.Height

	// Background
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), p.BackgroundColor, true)

	// Calculate fill width
	rangeVal := p.Max - p.Min
	if rangeVal > 0 {
		ratio := (p.Value - p.Min) / rangeVal
		fillW := w * ratio
		if fillW > 0 {
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(fillW), float32(h), p.FillColor, true)
		}
	}

	// Border
	// Top
	vector.StrokeLine(screen, float32(x), float32(y), float32(x+w), float32(y), 1, p.BorderColor, true)
	// Bottom
	vector.StrokeLine(screen, float32(x), float32(y+h), float32(x+w), float32(y+h), 1, p.BorderColor, true)
	// Left
	vector.StrokeLine(screen, float32(x), float32(y), float32(x), float32(y+h), 1, p.BorderColor, true)
	// Right
	vector.StrokeLine(screen, float32(x+w), float32(y), float32(x+w), float32(y+h), 1, p.BorderColor, true)
}
