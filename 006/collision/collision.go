package collision

import (
	"slices"

	"github.com/quasilyte/gmath"
)

// 衝突判定範囲のインターフェース
type Tester interface {
	Test(c Tester) bool
}

// 凸型多角形
type Polygon struct {
	Pos      gmath.Vec   // 座標
	Rad      gmath.Rad   // 回転角度(ラジアン)
	Vertices []gmath.Vec // 右周りの頂点集合
	Origin   gmath.Vec   // 頂点集合の回転原点
}

func (p *Polygon) Test(o Tester) bool {
	result := false
	switch v := o.(type) {
	case *Polygon:
		result = TestPolygonPolygon(v, p)
	case *Circle:
		result = TestCirclePolygon(v, p)
	case *Composit:
		result = TestPolygonComposit(p, v)
	default:
		result = false
	}

	return result
}

// 円
type Circle struct {
	Pos    gmath.Vec // 中心座標
	Radius float64   // 半径
}

func (c *Circle) Test(o Tester) bool {
	result := false
	switch v := o.(type) {
	case *Polygon:
		result = TestCirclePolygon(c, v)
	case *Circle:
		result = TestCircleCircle(c, v)
	case *Composit:
		result = TestCircleComposit(c, v)
	default:
		result = false
	}

	return result
}

// 複合形状
type Composit struct {
	Collisions []Tester         // 型的にCompositにCompositを入れることができそうだが禁止とする
	Operator   CompositOperator // 0:or、1:and
}

func (c *Composit) Test(o Tester) bool {
	result := false
	switch v := o.(type) {
	case *Polygon:
		result = TestPolygonComposit(v, c)
	case *Circle:
		result = TestCircleComposit(v, c)
	case *Composit:
		result = TestCompositComposit(v, c)
	default:
		result = false
	}

	return result
}

// 複合形状のAnd/Or条件
type CompositOperator int

const (
	CompositOr  CompositOperator = 0
	CompositAnd CompositOperator = 1
)

// 点と円の判定
func TestPointCircle(x, y float64, c1 *Circle) bool {
	return gmath.Vec{X: x, Y: y}.DistanceSquaredTo(c1.Pos) < c1.Radius*c1.Radius
}

// 点と凸型多角形の判定
func TestPointPolygon(x, y float64, c1 *Polygon) bool {
	// 外積を計算して境界内に点があるかを判定する
	r := make([]gmath.Vec, 0, len(c1.Vertices)+1)
	for _, p := range c1.Vertices {
		// 各頂点から回転原点を引いてから回転、回転原点とベース座標を足すことでグローバル座標を算出する
		v := p.Sub(c1.Origin).Rotated(c1.Rad).Add(c1.Origin).Add(c1.Pos)
		r = append(r, v)
	}

	// 1個目の頂点をスライスに追加する
	// これでr[n+1]-r[n]がエッジのベクトルになる
	r = append(r, r[0])

	// エッジのベクトルと、r[n]からx,yまでのベクトルの外積を取って、点がエッジのどちら側にあるかを調べる
	for i := 0; i < len(r)-1; i++ {
		v1 := gmath.Vec{X: x, Y: y}.Sub(r[i])
		v2 := r[i+1].Sub(r[i])

		// 外積
		cp := v1.X*v2.Y - v2.X*v1.Y

		// 外側にある
		if cp > 0 {
			// 当たってない
			return false
		}
	}

	// 全部内側にあれば当たっている
	return true
}

// 円同士の判定
func TestCircleCircle(c1 *Circle, c2 *Circle) bool {
	return c1.Pos.DistanceSquaredTo(c2.Pos) < (c1.Radius+c2.Radius)*(c1.Radius+c2.Radius)
}

// 円と凸型多角形の判定
func TestCirclePolygon(c *Circle, p *Polygon) bool {
	r := make([]gmath.Vec, 0, len(p.Vertices)+1)
	for _, v := range p.Vertices {
		// 各頂点から回転原点を引いてから回転、回転原点とベース座標を足すことでグローバル座標を算出する
		v := v.Sub(p.Origin).Rotated(p.Rad).Add(p.Origin).Add(p.Pos)
		r = append(r, v)
	}

	// 1個目の頂点をスライスに追加する
	// これでr1[n+1]-r1[n]がエッジのベクトルになる
	r = append(r, r[0])

	// エッジと円の中心点の位置関係を調べる
	norms := make([]gmath.Vec, len(p.Vertices))      // エッジの正規化ベクトル
	edges := make([]gmath.Vec, len(p.Vertices))      // エッジのベクトル
	points := make([]gmath.Vec, len(p.Vertices))     // エッジの起点から円の中心へのベクトル
	inside := true                                   // 円の中心が多角形の内側にあるフラグ
	before_edge_flg := make([]bool, len(p.Vertices)) // エッジより前に円の中心があるフラグ
	after_edge_flg := make([]bool, len(p.Vertices))  // エッジより後に円の中心があるフラグ

	for i := 0; i < len(r)-1; i++ {
		edges[i] = r[i+1].Sub(r[i])
		norms[i] = edges[i].Normalized()
		points[i] = c.Pos.Sub(r[i])

		// エッジの端より前後にあるかを内積で求める
		v := points[i].Dot(norms[i])
		e := edges[i].Dot(norms[i])
		before_edge_flg[i] = v <= 0 // 原点よりも手前
		after_edge_flg[i] = v >= e  // エッジの終点よりも先

		if i > 0 { // 最初の1個の場合は前のエッジが無い
			// 前のエッジよりも先で今のエッジよりも前の場合、エッジの起点の頂点が円に最も近い
			if after_edge_flg[i-1] && before_edge_flg[i] {
				// この頂点以外との衝突判定は必要ない
				return TestPointCircle(r[i].X, r[i].Y, c)
			}
		}

		// 外積を求めてどっち側にあるかを調べる
		cp := points[i].X*norms[i].Y - norms[i].X*points[i].Y
		if cp > 0 {
			// 外側にある
			inside = false

			// 頂点よりエッジのほうが近い場合、エッジと円の衝突判定をする
			if !before_edge_flg[i] && !after_edge_flg[i] {
				// エッジから円の中心までの最短距離は計算済
				if cp < c.Radius {
					return true
				}
			}
		}
	}

	// 最後のエッジよりも先で1個目のエッジよりも前の場合、1個目のエッジの起点の頂点が円に最も近い
	if after_edge_flg[len(r)-2] && before_edge_flg[0] {
		return TestPointCircle(r[0].X, r[0].Y, c)
	}

	// すべてのエッジの内側であれば衝突している
	return inside
}

// 凸型多角形同士の判定(SAT)
func TestPolygonPolygon(c1 *Polygon, c2 *Polygon) bool {
	r1 := make([]gmath.Vec, 0, len(c1.Vertices)+1)
	for _, p := range c1.Vertices {
		// 各頂点から回転原点を引いてから回転、回転原点とベース座標を足すことでグローバル座標を算出する
		v := p.Sub(c1.Origin).Rotated(c1.Rad).Add(c1.Origin).Add(c1.Pos)
		r1 = append(r1, v)
	}

	// 1個目の頂点をスライスに追加する
	// これでr1[n+1]-r1[n]がエッジのベクトルになる
	r1 = append(r1, r1[0])

	r2 := make([]gmath.Vec, 0, len(c2.Vertices)+1)
	for _, p := range c2.Vertices {
		// 各頂点から回転原点を引いてから回転、回転原点とベース座標を足すことでグローバル座標を算出する
		v := p.Sub(c2.Origin).Rotated(c2.Rad).Add(c2.Origin).Add(c2.Pos)
		r2 = append(r2, v)
	}

	// 1個目の頂点をスライスに追加する
	// これでr2[n+1]-r2[n]がエッジのベクトルになる
	r2 = append(r2, r2[0])

	// 軸ベクトルを算出
	vs := make([]gmath.Vec, 0, len(c1.Vertices)+len(c2.Vertices))
	for i := 0; i < len(c1.Vertices); i++ {
		vs = append(vs, r1[i+1].Sub(r1[i]).Normalized())
	}
	for i := 0; i < len(c2.Vertices); i++ {
		vs = append(vs, r2[i+1].Sub(r2[i]).Normalized())
	}

	// 各軸に各頂点を射影してmin/maxを求める
	for _, v := range vs {
		s1 := make([]float64, 0, 4)

		// c1の頂点を射影
		for i := 0; i < len(c1.Vertices); i++ {
			s1 = append(s1, r1[i].X*v.Y-v.X*r1[i].Y)
		}

		s2 := make([]float64, 0, 4)

		// c2の頂点を射影
		for i := 0; i < len(c2.Vertices); i++ {
			s2 = append(s2, r2[i].X*v.Y-v.X*r2[i].Y)
		}

		// min/max算出
		min1 := slices.Min(s1)
		max1 := slices.Max(s1)
		min2 := slices.Min(s2)
		max2 := slices.Max(s2)

		// 範囲が重なっているかのチェック
		if min1 > max2 || max1 < min2 {
			// 重なっていない
			return false
		}
	}

	// すべての軸で重なっていた
	return true
}

// 点と複合形状の判定
func TestPointComposit(x, y float64, co *Composit) bool {
	result := false
	for _, d := range co.Collisions {
		switch v := d.(type) {
		case *Polygon:
			result = TestPointPolygon(x, y, v)
		case *Circle:
			result = TestPointCircle(x, y, v)
		default:
			result = false
		}

		if co.Operator == CompositOr && result {
			return true
		}
		if co.Operator == CompositAnd && !result {
			return false
		}
	}

	return result
}

// 円と複合形状の判定
func TestCircleComposit(c *Circle, co *Composit) bool {
	result := false
	for _, d := range co.Collisions {
		result = c.Test(d)

		if co.Operator == CompositOr && result {
			return true
		}
		if co.Operator == CompositAnd && !result {
			return false
		}
	}

	return result
}

func TestPolygonComposit(p *Polygon, co *Composit) bool {
	result := false
	for _, d := range co.Collisions {
		result = p.Test(d)

		if co.Operator == CompositOr && result {
			return true
		}
		if co.Operator == CompositAnd && !result {
			return false
		}
	}

	return result
}

func TestCompositComposit(c1 *Composit, c2 *Composit) bool {
	result := false
	for _, d1 := range c1.Collisions {
		for _, d2 := range c2.Collisions {
			result = d1.Test(d2)

			if c2.Operator == CompositOr && result {
				break
			}
			if c2.Operator == CompositAnd && !result {
				break
			}
		}

		if c1.Operator == CompositOr && result {
			return true
		}
		if c1.Operator == CompositAnd && !result {
			return false
		}
	}

	return result
}
