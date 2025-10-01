package primitive

import (
	"image"
	"image/color"
	"math"

	"myproject/collision"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/gmath"
)

var emptyImage = ebiten.NewImage(3, 3)
var whitePixel = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

func init() {
	emptyImage.Fill(color.White)
}

// ちょっと乱雑になっていまっている
type Object interface {
	Draw(screen *ebiten.Image)
	Update()
	TestCollinsion(Object) bool
	Move(fx, fy, tx, ty float64)
	CheckPoint(x, y float64) bool
	SetFillColor(color.Color)
	GetComposit() *collision.Composit
	GetPos() gmath.Vec
}

type Base struct {
	Pos                gmath.Vec
	Rad                gmath.Rad
	FillColor          color.Color
	collision.Composit // 処理の簡素化のためにComposit専用とする
}

func NewPolygon(x, y, r float64, vs []gmath.Vec) *Base {
	c := collision.Composit{
		Collisions: []collision.Tester{
			&collision.Polygon{
				Pos:      gmath.Vec{X: x, Y: y},
				Rad:      gmath.Rad(r),
				Origin:   gmath.Vec{X: 0, Y: 0},
				Vertices: vs,
			},
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

// 特殊な形状を除いて、基本的には衝突判定の範囲を描画する
func (b *Base) Draw(screen *ebiten.Image) {
	for _, c := range b.Collisions {
		switch d := c.(type) {
		case *collision.Polygon: // 凸型多角形の描画
			var path vector.Path

			// 全座標の計算とパス設定
			pos := d.Vertices[0].Sub(d.Origin).Rotated(d.Rad).Add(d.Origin).Add(d.Pos)
			path.MoveTo(float32(pos.X), float32(pos.Y))

			for i := 1; i < len(d.Vertices); i++ {
				v := d.Vertices[i].Sub(d.Origin).Rotated(d.Rad).Add(d.Origin).Add(d.Pos)
				path.LineTo(float32(v.X), float32(v.Y))
			}
			path.Close()

			// 描画用頂点情報作成
			var vertices []ebiten.Vertex = []ebiten.Vertex{}
			var indices []uint16 = []uint16{}
			r, g, b, _ := b.FillColor.RGBA()
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
		case *collision.Circle: // 円の描画
			vector.DrawFilledCircle(screen, float32(d.Pos.X), float32(d.Pos.Y), float32(d.Radius), b.FillColor, true)
		}
	}
}

// 座標(x, y)がRectの中にあるかどうかをチェックする
func (b *Base) CheckPoint(x, y float64) bool {
	// 点と凸型多角形の衝突判定
	return collision.TestPointComposit(x, y, &b.Composit)
}

// ベクトル(x, y)を正規化する
func normalize(x, y float64) (float64, float64) {
	len := math.Hypot(x, y)
	if len != 0 {
		x /= len
		y /= len
	}
	return x, y
}

// 情報を更新する
func (b *Base) Update() {
	for _, c := range b.Collisions {
		switch d := c.(type) {
		case *collision.Polygon:
			d.Pos = b.Pos
			d.Rad = b.Rad
		case *collision.Circle:
			d.Pos = b.Pos
		}
	}
	b.FillColor = color.RGBA{0x00, 0xff, 0xff, 0xff}
}

// 重なっているかどうかをチェックする
func (b *Base) TestCollinsion(o Object) bool {
	return b.Test(o.GetComposit())
}

// マウスドラッグで掴んだ場所をfx,fyからtx,tyまで移動させる
func (b *Base) Move(fx, fy, tx, ty float64) {
	oldAngle := math.Atan2(b.Pos.Y-fy, b.Pos.X-fx)
	len := math.Hypot(fx-b.Pos.X, fy-b.Pos.Y)
	vx, vy := normalize(b.Pos.X-tx, b.Pos.Y-ty)
	b.Rad = b.Rad + gmath.Rad(math.Atan2(vy, vx)-oldAngle)
	b.Pos.X = tx + vx*len
	b.Pos.Y = ty + vy*len
}

func (b *Base) SetFillColor(c color.Color) {
	b.FillColor = c
}

func (b *Base) GetComposit() *collision.Composit {
	return &b.Composit
}

func (b *Base) GetPos() gmath.Vec {
	return b.Pos
}
