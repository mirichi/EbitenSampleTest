package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var emptyImage = ebiten.NewImage(3, 3)
var whitePixel = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

func init() {
	emptyImage.Fill(color.White)
}

type Rect struct {
	X, Y          float64
	Width, Height float64
	Rad           float64
	FillColor     color.RGBA
}

func NewRect(x, y, w, h float64) *Rect {
	return &Rect{
		X:         x,
		Y:         y,
		Width:     w,
		Height:    h,
		Rad:       math.Pi / 4,
		FillColor: color.RGBA{0x00, 0xff, 0xff, 0xff},
	}
}

func (r *Rect) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-0.5, -0.5)
	op.GeoM.Scale(r.Width, r.Height)
	op.GeoM.Rotate(r.Rad)
	op.GeoM.Translate(r.X, r.Y)
	op.ColorScale.ScaleWithColor(r.FillColor)
	screen.DrawImage(whitePixel, op)
}

// 座標(x, y)がRectの中にあるかどうかをチェックする
func (r *Rect) CheckPoint(x, y float64) bool {
	// 四角形の中心座標
	cx := r.X
	cy := r.Y

	// 座標(x, y)から中心へのベクトル
	dx := x - cx
	dy := y - cy

	// 回転を逆に戻す
	cos := math.Cos(r.Rad)
	sin := math.Sin(r.Rad)

	px := dx*cos + dy*sin
	py := -dx*sin + dy*cos

	hw := r.Width * 0.5
	hh := r.Height * 0.5
	return px >= -hw && px <= hw && py >= -hh && py <= hh
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

// Rect同士が重なっているかどうかをチェックする
func (r *Rect) CheckBox(r2 *Rect) bool {
	// rの4点の回転後座標を計算する
	points := [4][2]float64{}
	cos := math.Cos(r.Rad)
	sin := math.Sin(r.Rad)
	hw := r.Width * 0.5
	hh := r.Height * 0.5
	points[0][0] = r.X + (-hw*cos - -hh*sin)
	points[0][1] = r.Y + (-hw*sin + -hh*cos)
	points[1][0] = r.X + (hw*cos - -hh*sin)
	points[1][1] = r.Y + (hw*sin + -hh*cos)
	points[2][0] = r.X + (hw*cos - hh*sin)
	points[2][1] = r.Y + (hw*sin + hh*cos)
	points[3][0] = r.X + (-hw*cos - hh*sin)
	points[3][1] = r.Y + (-hw*sin + hh*cos)

	// r2の4点の回転後座標を計算する
	points2 := [4][2]float64{}
	cos2 := math.Cos(r2.Rad)
	sin2 := math.Sin(r2.Rad)
	hw2 := r2.Width * 0.5
	hh2 := r2.Height * 0.5
	points2[0][0] = r2.X + (-hw2*cos2 - -hh2*sin2)
	points2[0][1] = r2.Y + (-hw2*sin2 + -hh2*cos2)
	points2[1][0] = r2.X + (hw2*cos2 - -hh2*sin2)
	points2[1][1] = r2.Y + (hw2*sin2 + -hh2*cos2)
	points2[2][0] = r2.X + (hw2*cos2 - hh2*sin2)
	points2[2][1] = r2.Y + (hw2*sin2 + hh2*cos2)
	points2[3][0] = r2.X + (-hw2*cos2 - hh2*sin2)
	points2[3][1] = r2.Y + (-hw2*sin2 + hh2*cos2)

	// SATによる当たり判定
	axes := [4][2]float64{
		{points[1][0] - points[0][0], points[1][1] - points[0][1]},
		{points[3][0] - points[0][0], points[3][1] - points[0][1]},
		{points2[1][0] - points2[0][0], points2[1][1] - points2[0][1]},
		{points2[3][0] - points2[0][0], points2[3][1] - points2[0][1]},
	}
	for _, axis := range axes {
		// 軸を正規化
		axisX, axisY := normalize(axis[0], axis[1])
		// rの4点を軸に投影
		minA, maxA := math.Inf(1), math.Inf(-1)
		for _, p := range points {
			proj := p[0]*axisX + p[1]*axisY
			if proj < minA {
				minA = proj
			}
			if proj > maxA {
				maxA = proj
			}
		}
		// r2の4点を軸に投影
		minB, maxB := math.Inf(1), math.Inf(-1)
		for _, p := range points2 {
			proj := p[0]*axisX + p[1]*axisY
			if proj < minB {
				minB = proj
			}
			if proj > maxB {
				maxB = proj
			}
		}
		// 投影が重なっているかチェック
		if maxA < minB || maxB < minA {
			// 重なっていない場合はfalseを返す
			return false
		}
	}

	// すべての軸で重なっている場合はtrueを返す
	return true
}

// マウスドラッグで掴んだ場所をfx,fyからtx,tyまで移動させる
func (r *Rect) Move(fx, fy, tx, ty float64) {
	oldAngle := math.Atan2(r.Y-fy, r.X-fx)
	len := math.Hypot(fx-r.X, fy-r.Y)
	vx, vy := normalize(r.X-tx, r.Y-ty)
	r.Rad += math.Atan2(vy, vx) - oldAngle
	r.X = tx + vx*len
	r.Y = ty + vy*len
}
