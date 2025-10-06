package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jakecoffman/cp/v2"
)

var (
	space *cp.Space
	body  *cp.Body
	shape *cp.Shape
)

type Game struct{}

func addBall(space *cp.Space, pos cp.Vector, radius float64) {
	// 質量
	mass := 100.0

	// 円のモーメントを計算
	moment := cp.MomentForCircle(mass, 0, radius, cp.Vector{})

	// Bodyを作成
	body = cp.NewBody(mass, moment)

	// Bodyの位置を設定
	body.SetPosition(pos)

	// Shape(Circle)を作成
	shape = cp.NewCircle(body, radius, cp.Vector{})

	// 弾性を設定
	shape.SetElasticity(0.5)

	// 摩擦を設定
	shape.SetFriction(0.5)

	// Spaceに追加
	space.AddBody(body)
	space.AddShape(shape)
}

func addWall(space *cp.Space, p1, p2 cp.Vector, radius float64) {
	// Shape(Segment)を作成
	segment := cp.NewSegment(space.StaticBody, p1, p2, radius)

	// 弾性を設定
	segment.SetElasticity(0.5)

	// 摩擦を設定
	segment.SetFriction(0.5)

	// Spaceに追加
	space.AddShape(segment)
}

func (g *Game) Update() error {
	// 1/60秒分だけ動かす
	space.Step(1 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// bodyの位置に円を描画する
	vector.StrokeCircle(screen, float32(body.Position().X), float32(body.Position().Y), float32(50), 5, color.White, true)

	// wallの位置に線を描画する(変数に入れてないので固定で描く)
	vector.StrokeLine(screen, 40, 350, 400, 460, 5, color.White, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	space = cp.NewSpace()

	// 重力を設定
	space.SetGravity(cp.Vector{X: 0, Y: 100})

	// ボールを作成してSpaceに追加する
	addBall(space, cp.Vector{X: 200, Y: 50}, 50)

	// 線を作成してSpaceに追加する
	addWall(space, cp.Vector{X: 40, Y: 350}, cp.Vector{X: 400, Y: 460}, 5)

	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
