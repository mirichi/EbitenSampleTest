// Package collision provides basic collision detection functionalities
// for convex polygons and circles, including composite shapes with AND/OR logic.
// It uses the Separating Axis Theorem (SAT) for polygon-polygon collision detection
// and handles various combinations of shapes.
// 衝突判定機能を提供する
// 凸型多角形と円、複合形状のAND/OR条件をサポート
// 凸型多角形同士の衝突判定には分離軸定理(SAT)を使用
// 円と多角形、円同士、多角形と複合形状など様々な組み合わせに対応

package objects

import (
	"slices"
)

// 衝突判定範囲のインターフェース
type CollisionTester interface {
	Test(c CollisionTester) bool
	CollisionShape() CollisionTester // 実体の衝突判定オブジェクトを返す
}

// 凸型多角形
type PolygonCollider struct {
	Sprite   *Sprite
	Vertices []Vector2 // 右周りの頂点集合
}

// NewPolygonCollider は新しい凸型多角形の衝突判定を作成する
// vertices: Spriteの画像左上を(0,0)としたローカル座標系で指定してください
func NewPolygonCollider(sprite *Sprite, vertices []Vector2) *PolygonCollider {
	return &PolygonCollider{Sprite: sprite, Vertices: vertices}
}

// NewRectCollider はSpriteの画像サイズに合わせて四角形の衝突判定を作成します
func NewRectCollider(sprite *Sprite) *PolygonCollider {
	bounds := sprite.Image.Bounds()
	w, h := float64(bounds.Dx()), float64(bounds.Dy())
	vertices := []Vector2{
		{0, 0},
		{w, 0},
		{w, h},
		{0, h},
	}
	return NewPolygonCollider(sprite, vertices)
}

func (p *PolygonCollider) Test(o CollisionTester) bool {
	o = o.CollisionShape() // 実体を取り出す
	result := false
	switch v := o.(type) {
	case *PolygonCollider:
		result = TestPolygonPolygon(v, p)
	case *CircleCollider:
		result = TestCirclePolygon(v, p)
	case *CompositeCollider:
		result = TestPolygonComposit(p, v)
	default:
		result = false
	}

	return result
}

func (p *PolygonCollider) CollisionShape() CollisionTester {
	return p
}

// 円
type CircleCollider struct {
	Sprite *Sprite
	Center Vector2
	Radius float64 // 半径
}

// NewCircleCollider は新しい円の衝突判定を作成する
// pos: Spriteの画像左上を(0,0)としたローカル座標系で指定してください
func NewCircleCollider(sprite *Sprite, pos Vector2, radius float64) *CircleCollider {
	return &CircleCollider{Sprite: sprite, Center: pos, Radius: radius}
}

func (c *CircleCollider) Test(o CollisionTester) bool {
	o = o.CollisionShape() // 実体を取り出す
	result := false
	switch v := o.(type) {
	case *PolygonCollider:
		result = TestCirclePolygon(c, v)
	case *CircleCollider:
		result = TestCircleCircle(c, v)
	case *CompositeCollider:
		result = TestCircleComposit(c, v)
	default:
		result = false
	}

	return result
}

func (c *CircleCollider) CollisionShape() CollisionTester {
	return c
}

// 複合形状
type CompositeCollider struct {
	Collisions []CollisionTester
	Operator   CompositeOperator // 0:or、1:and
}

func NewCompositeCollider(collisions []CollisionTester, operator CompositeOperator) *CompositeCollider {
	return &CompositeCollider{Collisions: collisions, Operator: operator}
}

func (c *CompositeCollider) Test(o CollisionTester) bool {
	o = o.CollisionShape() // 実体を取り出す
	result := false
	switch v := o.(type) {
	case *PolygonCollider:
		result = TestPolygonComposit(v, c)
	case *CircleCollider:
		result = TestCircleComposit(v, c)
	case *CompositeCollider:
		result = TestCompositComposit(v, c)
	default:
		result = false
	}

	return result
}

func (c *CompositeCollider) CollisionShape() CollisionTester {
	return c
}

// 複合形状のAnd/Or条件
type CompositeOperator int

const (
	CompositeOr  CompositeOperator = 0
	CompositeAnd CompositeOperator = 1
)

// 点と円の判定
func TestPointCircle(pos Vector2, c1 *CircleCollider) bool {
	// グローバル行列を使用して円の中心座標を変換
	matrix := c1.Sprite.GetGlobalMatrix()
	gx, gy := matrix.Apply(c1.Center.X, c1.Center.Y)

	cpos := Vector2{X: gx, Y: gy}
	return pos.DistanceSquaredTo(cpos) < c1.Radius*c1.Radius
}

// 点と凸型多角形の判定
func TestPointPolygon(pos Vector2, p1 *PolygonCollider) bool {
	// グローバル行列を取得
	matrix := p1.Sprite.GetGlobalMatrix()

	r := make([]Vector2, 0, len(p1.Vertices)+1)
	for _, p := range p1.Vertices {
		gx, gy := matrix.Apply(p.X, p.Y)
		r = append(r, Vector2{X: gx, Y: gy})
	}

	// 1個目の頂点をスライスに追加する
	r = append(r, r[0])

	// エッジのベクトルと、r[n]からx,yまでのベクトルの外積を取って、点がエッジのどちら側にあるかを調べる
	for i := 0; i < len(r)-1; i++ {
		v1 := pos.Sub(r[i])
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
func TestCircleCircle(c1 *CircleCollider, c2 *CircleCollider) bool {
	m1 := c1.Sprite.GetGlobalMatrix()
	gx1, gy1 := m1.Apply(c1.Center.X, c1.Center.Y)
	cpos1 := Vector2{X: gx1, Y: gy1}

	m2 := c2.Sprite.GetGlobalMatrix()
	gx2, gy2 := m2.Apply(c2.Center.X, c2.Center.Y)
	cpos2 := Vector2{X: gx2, Y: gy2}

	return cpos1.DistanceSquaredTo(cpos2) < (c1.Radius+c2.Radius)*(c1.Radius+c2.Radius)
}

// 円と凸型多角形の判定
func TestCirclePolygon(c1 *CircleCollider, p1 *PolygonCollider) bool {
	m1 := c1.Sprite.GetGlobalMatrix()
	gx1, gy1 := m1.Apply(c1.Center.X, c1.Center.Y)
	cpos := Vector2{X: gx1, Y: gy1}

	m2 := p1.Sprite.GetGlobalMatrix()

	r := make([]Vector2, 0, len(p1.Vertices)+1)
	for _, v := range p1.Vertices {
		gx, gy := m2.Apply(v.X, v.Y)
		r = append(r, Vector2{X: gx, Y: gy})
	}

	// 1個目の頂点をスライスに追加する
	// これでr1[n+1]-r1[n]がエッジのベクトルになる
	r = append(r, r[0])

	// エッジと円の中心点の位置関係を調べる
	norms := make([]Vector2, len(p1.Vertices))        // エッジの正規化ベクトル
	edges := make([]Vector2, len(p1.Vertices))        // エッジのベクトル
	points := make([]Vector2, len(p1.Vertices))       // エッジの起点から円の中心へのベクトル
	inside := true                                    // 円の中心が多角形の内側にあるフラグ
	before_edge_flg := make([]bool, len(p1.Vertices)) // エッジより前に円の中心があるフラグ
	after_edge_flg := make([]bool, len(p1.Vertices))  // エッジより後に円の中心があるフラグ

	// 中略なしで実装を戻す必要があるが、replace_file_contentの制限でここだけ修正するとしても周囲のコンテキストが必要。
	// ここではTestCirclePolygonの後半は変更がないため、前半の座標計算部分とループ初期化のみを置き換える形にするのが安全だが、
	// replace_file_contentはブロック置換なので、関数全体を再定義する。

	for i := 0; i < len(r)-1; i++ {
		edges[i] = r[i+1].Sub(r[i])
		norms[i] = edges[i].Normalize()
		points[i] = cpos.Sub(r[i])

		// エッジの端より前後にあるかを内積で求める
		v := points[i].Dot(norms[i])
		e := edges[i].Dot(norms[i])
		before_edge_flg[i] = v <= 0 // 原点よりも手前
		after_edge_flg[i] = v >= e  // エッジの終点よりも先

		if i > 0 { // 最初の1個の場合は前のエッジが無い
			// 前のエッジよりも先で今のエッジよりも前の場合、エッジの起点の頂点が円に最も近い
			if after_edge_flg[i-1] && before_edge_flg[i] {
				// この頂点以外との衝突判定は必要ない
				// return TestPointCircle(r[i], c)
				return r[i].DistanceSquaredTo(cpos) < c1.Radius*c1.Radius
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
				if cp < c1.Radius {
					return true
				}
			}
		}
	}

	// 最後のエッジよりも先で1個目のエッジよりも前の場合、1個目のエッジの起点の頂点が円に最も近い
	if after_edge_flg[len(r)-2] && before_edge_flg[0] {
		return r[0].DistanceSquaredTo(cpos) < c1.Radius*c1.Radius
	}

	// すべてのエッジの内側であれば衝突している
	return inside
}

// 凸型多角形同士の判定(SAT)
func TestPolygonPolygon(c1 *PolygonCollider, c2 *PolygonCollider) bool {
	m1 := c1.Sprite.GetGlobalMatrix()
	r1 := make([]Vector2, 0, len(c1.Vertices)+1)
	for _, p := range c1.Vertices {
		gx, gy := m1.Apply(p.X, p.Y)
		r1 = append(r1, Vector2{X: gx, Y: gy})
	}

	// 1個目の頂点をスライスに追加する
	r1 = append(r1, r1[0])

	m2 := c2.Sprite.GetGlobalMatrix()
	r2 := make([]Vector2, 0, len(c2.Vertices)+1)
	for _, p := range c2.Vertices {
		gx, gy := m2.Apply(p.X, p.Y)
		r2 = append(r2, Vector2{X: gx, Y: gy})
	}

	// 1個目の頂点をスライスに追加する
	r2 = append(r2, r2[0])

	// 軸ベクトルを算出
	vs := make([]Vector2, 0, len(c1.Vertices)+len(c2.Vertices))
	for i := 0; i < len(c1.Vertices); i++ {
		vs = append(vs, r1[i+1].Sub(r1[i]).Normalize())
	}
	for i := 0; i < len(c2.Vertices); i++ {
		vs = append(vs, r2[i+1].Sub(r2[i]).Normalize())
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
func TestPointComposit(pos Vector2, co *CompositeCollider) bool {
	result := false
	for _, d := range co.Collisions {
		switch v := d.(type) {
		case *PolygonCollider:
			result = TestPointPolygon(pos, v)
		case *CircleCollider:
			result = TestPointCircle(pos, v)
		default:
			result = false
		}

		if co.Operator == CompositeOr && result {
			return true
		}
		if co.Operator == CompositeAnd && !result {
			return false
		}
	}

	return result
}

// 円と複合形状の判定
func TestCircleComposit(c *CircleCollider, co *CompositeCollider) bool {
	result := false
	for _, d := range co.Collisions {
		result = c.Test(d)

		if co.Operator == CompositeOr && result {
			return true
		}
		if co.Operator == CompositeAnd && !result {
			return false
		}
	}

	return result
}

func TestPolygonComposit(p *PolygonCollider, co *CompositeCollider) bool {
	result := false
	for _, d := range co.Collisions {
		result = p.Test(d)

		if co.Operator == CompositeOr && result {
			return true
		}
		if co.Operator == CompositeAnd && !result {
			return false
		}
	}

	return result
}

func TestCompositComposit(c1 *CompositeCollider, c2 *CompositeCollider) bool {
	result := false
	for _, d1 := range c1.Collisions {
		for _, d2 := range c2.Collisions {
			result = d1.Test(d2)

			if c2.Operator == CompositeOr && result {
				break
			}
			if c2.Operator == CompositeAnd && !result {
				break
			}
		}

		if c1.Operator == CompositeOr && result {
			return true
		}
		if c1.Operator == CompositeAnd && !result {
			return false
		}
	}

	return result
}
