package main

import (
	"image"
	"image/color"
	"math"
	"myproject/control"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var emptyImage = ebiten.NewImage(3, 3)
var whitePixel = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

func init() {
	emptyImage.Fill(color.White)
}

// 3次ベジェ曲線座標算出
func (g *Game) tesselate_bezier(x1, y1, x2, y2, x3, y3, x4, y4, level float64, p *vector.Path) {
	// 10回までしか再帰しない
	if level > 10 {
		return
	}

	dx := x4 - x1
	dy := y4 - y1
	d2 := math.Abs((x2-x4)*dy - (y2-y4)*dx)
	d3 := math.Abs((x3-x4)*dy - (y3-y4)*dx)

	// この範囲が直線なら終了
	if (d2+d3)*(d2+d3) < 0.25*(dx*dx+dy*dy) {
		p.LineTo(float32(x4), float32(y4))
		return
	}

	// 2分割
	x12 := (x1 + x2) * 0.5
	y12 := (y1 + y2) * 0.5
	x23 := (x2 + x3) * 0.5
	y23 := (y2 + y3) * 0.5
	x34 := (x3 + x4) * 0.5
	y34 := (y3 + y4) * 0.5
	x123 := (x12 + x23) * 0.5
	y123 := (y12 + y23) * 0.5
	x234 := (x23 + x34) * 0.5
	y234 := (y23 + y34) * 0.5
	x1234 := (x123 + x234) * 0.5
	y1234 := (y123 + y234) * 0.5

	// 分割した前半分を処理
	g.tesselate_bezier(x1, y1, x12, y12, x123, y123, x1234, y1234, level+1, p)

	// 後ろ半分を処理
	g.tesselate_bezier(x1234, y1234, x234, y234, x34, y34, x4, y4, level+1, p)
}

// 3次ベジェ曲線描画
func (g *Game) draw_bezier(screen *ebiten.Image) {
	var path vector.Path

	p0 := g.objects[0].GetPos() // 始点
	p1 := g.objects[1].GetPos() // 制御点1
	p2 := g.objects[2].GetPos() // 制御点2
	p3 := g.objects[3].GetPos() // 終点

	path.MoveTo(float32(p0.X), float32(p0.Y))
	g.tesselate_bezier(p0.X, p0.Y, p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, 0, &path) // pathの前にある0は再帰の深さなので0固定で呼ぶ

	// Strokeで描画
	op := &vector.StrokeOptions{}
	op.Width = 10
	op.LineJoin = vector.LineJoinRound
	op.LineCap = vector.LineCapRound
	var vertices []ebiten.Vertex = []ebiten.Vertex{}
	var indices []uint16 = []uint16{}
	vertices, indices = path.AppendVerticesAndIndicesForStroke(vertices[:0], indices[:0], op)
	for i := range vertices {
		vertices[i].SrcX = 1
		vertices[i].SrcY = 1
		vertices[i].ColorR = 0xbf / float32(0xff)
		vertices[i].ColorG = 0xbf / float32(0xff)
		vertices[i].ColorB = 0x30 / float32(0xff)
		vertices[i].ColorA = 1
	}
	op2 := &ebiten.DrawTrianglesOptions{}
	op2.AntiAlias = true
	op2.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(vertices, indices, whitePixel, op2)

	// 制御点をつなぐ線を描画する
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}
	green := color.RGBA{0x00, 0xff, 0x00, 0xff}
	red := color.RGBA{0xff, 0x00, 0x00, 0xff}
	s := g.controls[0].(*control.Slider)
	p := float64(s.X) / (640 - float64(s.W))

	vector.StrokeLine(screen, float32(p0.X), float32(p0.Y), float32(p1.X), float32(p1.Y), 1, white, true)
	vector.StrokeLine(screen, float32(p1.X), float32(p1.Y), float32(p2.X), float32(p2.Y), 1, white, true)
	vector.StrokeLine(screen, float32(p2.X), float32(p2.Y), float32(p3.X), float32(p3.Y), 1, white, true)

	pp0 := p0.Add(p1.Sub(p0).Mulf(p))
	pp1 := p1.Add(p2.Sub(p1).Mulf(p))
	pp2 := p2.Add(p3.Sub(p2).Mulf(p))

	vector.DrawFilledCircle(screen, float32(pp0.X), float32(pp0.Y), float32(3), white, true)
	vector.DrawFilledCircle(screen, float32(pp1.X), float32(pp1.Y), float32(3), white, true)
	vector.DrawFilledCircle(screen, float32(pp2.X), float32(pp2.Y), float32(3), white, true)

	vector.StrokeLine(screen, float32(pp0.X), float32(pp0.Y), float32(pp1.X), float32(pp1.Y), 1, green, true)
	vector.StrokeLine(screen, float32(pp1.X), float32(pp1.Y), float32(pp2.X), float32(pp2.Y), 1, green, true)

	ppp0 := pp0.Add(pp1.Sub(pp0).Mulf(p))
	ppp1 := pp1.Add(pp2.Sub(pp1).Mulf(p))

	vector.DrawFilledCircle(screen, float32(ppp0.X), float32(ppp0.Y), float32(3), green, true)
	vector.DrawFilledCircle(screen, float32(ppp1.X), float32(ppp1.Y), float32(3), green, true)

	vector.StrokeLine(screen, float32(ppp0.X), float32(ppp0.Y), float32(ppp1.X), float32(ppp1.Y), 1, red, true)

	pppp := ppp0.Add(ppp1.Sub(ppp0).Mulf(p))
	vector.DrawFilledCircle(screen, float32(pppp.X), float32(pppp.Y), float32(3), red, true)
}
