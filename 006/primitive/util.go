package primitive

import (
	"image/color"
	"math"

	"myproject/collision"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/gmath"
)

func NewRect(x, y, w, h, r float64) *Base {
	c := collision.Polygon{
		Origin: gmath.Vec{X: 0, Y: 0},
		Vertices: []gmath.Vec{
			{X: -w / 2, Y: -h / 2},
			{X: w / 2, Y: -h / 2},
			{X: w / 2, Y: h / 2},
			{X: -w / 2, Y: h / 2},
		},
	}

	return &Base{
		Pos:       gmath.Vec{X: x, Y: y},
		Rad:       gmath.Rad(r),
		FillColor: color.RGBA{0x00, 0xff, 0xff, 0xff},
		Composit: collision.Composit{
			Collisions: []collision.Tester{&c},
		},
	}
}

func NewCircle(x, y, r float64) *Base {
	c := collision.Circle{
		Pos:    gmath.Vec{X: x, Y: y},
		Radius: r,
	}

	return &Base{
		Pos:       gmath.Vec{X: x, Y: y},
		FillColor: color.RGBA{0x00, 0xff, 0xff, 0xff},
		Composit: collision.Composit{
			Collisions: []collision.Tester{&c},
		},
	}
}

func NewStar(x, y, w, r float64) *Base {
	v1 := []gmath.Vec{}
	for i := 0; i < 5; i++ {
		v1 = append(v1, gmath.Vec{X: 0, Y: -1}.Rotated(gmath.DegToRad(float64(i*72))).Mulf(w))
	}

	v2 := []gmath.Vec{}
	for i := 0; i < 5; i++ {
		v2 = append(v2, gmath.Vec{X: 0, Y: -0.382}.Rotated(gmath.DegToRad(float64(i*72+36))).Mulf(w))
	}

	c := collision.Composit{
		Collisions: []collision.Tester{
			&collision.Polygon{
				Origin: gmath.Vec{X: 0, Y: 0},
				Vertices: []gmath.Vec{
					{X: v1[0].X, Y: v1[0].Y},
					{X: v2[1].X, Y: v2[1].Y},
					{X: v2[3].X, Y: v2[3].Y},
				}},
			&collision.Polygon{
				Origin: gmath.Vec{X: 0, Y: 0},
				Vertices: []gmath.Vec{
					{X: v1[1].X, Y: v1[1].Y},
					{X: v2[2].X, Y: v2[2].Y},
					{X: v2[4].X, Y: v2[4].Y},
				}},
			&collision.Polygon{
				Origin: gmath.Vec{X: 0, Y: 0},
				Vertices: []gmath.Vec{
					{X: v1[2].X, Y: v1[2].Y},
					{X: v2[3].X, Y: v2[3].Y},
					{X: v2[0].X, Y: v2[0].Y},
				}},
			&collision.Polygon{
				Origin: gmath.Vec{X: 0, Y: 0},
				Vertices: []gmath.Vec{
					{X: v1[3].X, Y: v1[3].Y},
					{X: v2[4].X, Y: v2[4].Y},
					{X: v2[1].X, Y: v2[1].Y},
				}},
			&collision.Polygon{
				Origin: gmath.Vec{X: 0, Y: 0},
				Vertices: []gmath.Vec{
					{X: v1[4].X, Y: v1[4].Y},
					{X: v2[0].X, Y: v2[0].Y},
					{X: v2[2].X, Y: v2[2].Y},
				}},
		},
		Operator: collision.CompositOr,
	}
	return &Base{
		Pos:       gmath.Vec{X: x, Y: y},
		Rad:       gmath.Rad(r),
		FillColor: color.RGBA{0x00, 0xff, 0xff, 0xff},
		Composit:  c,
	}
}

// and条件の確認用の半円
type HarfCircle struct {
	Base
	Radius float64
}

func NewHarfCircle(x, y, r float64) *HarfCircle {
	c := []collision.Tester{
		&collision.Circle{
			Pos:    gmath.Vec{X: x, Y: y},
			Radius: r,
		},
		&collision.Polygon{
			Origin: gmath.Vec{X: 0, Y: 0},
			Vertices: []gmath.Vec{
				{X: 0, Y: -r},
				{X: r, Y: -r},
				{X: r, Y: r},
				{X: 0, Y: r},
			},
		},
	}

	return &HarfCircle{
		Base: Base{
			Pos:       gmath.Vec{X: x, Y: y},
			Rad:       gmath.Rad(r),
			FillColor: color.RGBA{0x00, 0xff, 0xff, 0xff},
			Composit: collision.Composit{
				Collisions: c,
				Operator:   collision.CompositAnd,
			},
		},
		Radius: r,
	}
}

// HarfCircleは特殊な形状なので自前で描画する
func (c *HarfCircle) Draw(screen *ebiten.Image) {
	var path vector.Path

	// 半円描画
	path.MoveTo(float32(c.Pos.X), float32(c.Pos.Y))
	path.Arc(float32(c.Pos.X), float32(c.Pos.Y), float32(c.Radius), float32(c.Rad)-math.Pi*0.5, float32(c.Rad)+math.Pi*0.5, vector.Clockwise)
	path.Close()

	// 描画用頂点情報作成
	var vertices []ebiten.Vertex = []ebiten.Vertex{}
	var indices []uint16 = []uint16{}
	r, g, b, _ := c.FillColor.RGBA()
	vertices, indices = path.AppendVerticesAndIndicesForFilling(vertices[:0], indices[:0])
	for i := range vertices {
		vertices[i].SrcX = 1
		vertices[i].SrcY = 1
		vertices[i].ColorR = float32(r) / float32(0xff)
		vertices[i].ColorG = float32(g) / float32(0xff)
		vertices[i].ColorB = float32(b) / float32(0xff)
		vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(vertices, indices, whitePixel, op)
}
