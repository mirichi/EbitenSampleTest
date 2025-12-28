package parts

import "github.com/hajimehoshi/ebiten/v2"

var nextCursorShape = ebiten.CursorShapeDefault

// RequestCursorShape requests a cursor shape change for the next frame.
// If multiple requests are made, the last non-default one wins (or simply the last one).
func RequestCursorShape(shape ebiten.CursorShapeType) {
	if shape != ebiten.CursorShapeDefault {
		nextCursorShape = shape
	}
}

// FinalizeCursor applies the requested cursor shape and resets the buffer.
// This should be called once per frame, preferably at the end of Update.
func FinalizeCursor() {
	ebiten.SetCursorShape(nextCursorShape)
	nextCursorShape = ebiten.CursorShapeDefault
}
