package objects

import (
	"math"
	"testing"
)

// Helper to create a dummy sprite
func newTestSprite(x, y float64) *Sprite {
	return NewSprite(x, y, nil)
}

func TestPointCircleLogic(t *testing.T) {
	s := newTestSprite(100, 100)
	c := NewCircleCollision(s, Vector2{0, 0}, 10) // Radius 10, Center relative to sprite is (0,0)

	tests := []struct {
		name string
		pos  Vector2
		want bool
	}{
		{"Inside Center", Vector2{100, 100}, true},
		{"Inside Edge", Vector2{105, 100}, true},
		{"Outside", Vector2{111, 100}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TestPointCircle(tt.pos, c); got != tt.want {
				t.Errorf("TestPointCircle(%v) = %v, want %v", tt.pos, got, tt.want)
			}
		})
	}
}

func TestCircleCircleLogic(t *testing.T) {
	// Sprite at (100,100)
	s1 := newTestSprite(100, 100)
	// Circle center at (0,0) relative to sprite origin (which is 0,0 default)
	// So global center is (100,100)
	c1 := NewCircleCollision(s1, Vector2{0, 0}, 10)

	// Sprite at (115,100)
	s2 := newTestSprite(115, 100)
	c2 := NewCircleCollision(s2, Vector2{0, 0}, 10)

	// 100->115 dist is 15. Rad sum is 20. Collide.
	if !TestCircleCircle(c1, c2) {
		gx1, gy1 := s1.GetGlobalPos()
		gx2, gy2 := s2.GetGlobalPos()
		t.Errorf("Expected collision between circle at (%f,%f) and (%f,%f)", gx1, gy1, gx2, gy2)
	}

	// Move s2 to (121, 100). Dist 21. No collide.
	s2.X = 121
	if TestCircleCircle(c1, c2) {
		t.Errorf("Expected NO collision at dist 21")
	}
}

func TestPolygonPolygonLogic_Rotate(t *testing.T) {
	// Square 20x20 centered at local (0,0)
	// Global pos will be sprite pos
	vertices := []Vector2{
		{-10, -10}, {10, -10}, {10, 10}, {-10, 10},
	}

	s1 := newTestSprite(100, 100)
	p1 := NewPolygonCollision(s1, vertices)

	s2 := newTestSprite(125, 100) // Gap is 5 (110 vs 115) if aligned
	p2 := NewPolygonCollision(s2, vertices)

	// 100+10=110 (right edge of s1)
	// 125-10=115 (left edge of s2)
	// Gap = 5. No collision.
	if TestPolygonPolygon(p1, p2) {
		t.Error("Squares should not collide yet")
	}

	// Rotate s1 by 45 degrees.
	// Corner distance will extend approx 1.414 * 10 = 14.14
	// Rightmost point ~= 100 + 14.14 = 114.14
	// s2 leftmost is 115. Still no collision (very close).
	s1.Angle = math.Pi / 4
	if TestPolygonPolygon(p1, p2) {
		t.Error("Squares should not collide after 45deg rotation")
	}

	// Move s2 closer: 124
	// Leftmost: 114.
	// 114 < 114.14 -> Should collide
	s2.X = 124
	if !TestPolygonPolygon(p1, p2) {
		t.Error("Squares SHOULD collide after moving closer with rotation")
	}
}

func TestCirclePolygonLogic(t *testing.T) {
	// Square 20x20
	vertices := []Vector2{
		{-10, -10}, {10, -10}, {10, 10}, {-10, 10},
	}
	sPoly := newTestSprite(100, 100)
	p := NewPolygonCollision(sPoly, vertices)

	// Circle radius 5
	sCircle := newTestSprite(116, 100)
	c := NewCircleCollision(sCircle, Vector2{0, 0}, 5)

	// Poly Edge X=110. Circle X=116, R=5 -> Left=111.
	// Gap 1. No collision.
	if c.Test(p) { // c.Test(p) internally calls TestCirclePolygon
		t.Error("Circle and Polygon should NOT collide")
	}

	// Move Circle to 114. Left=109. Overlaps 110.
	sCircle.X = 114
	if !c.Test(p) {
		t.Error("Circle and Polygon SHOULD collide")
	}
}
