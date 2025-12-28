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
}

// 凸型多角形
type PolygonCollision struct {
	Sprite   *Sprite
	Vertices []Vector2 // 右周りの頂点集合
}

func NewPolygonCollision(sprite *Sprite, vertices []Vector2) *PolygonCollision {
	return &PolygonCollision{Sprite: sprite, Vertices: vertices}
}

func (p *PolygonCollision) Test(o CollisionTester) bool {
	result := false
	switch v := o.(type) {
	case *PolygonCollision:
		result = TestPolygonPolygon(v, p)
	case *CircleCollision:
		result = TestCirclePolygon(v, p)
	case *CompositeCollision:
		result = TestPolygonComposit(p, v)
	default:
		result = false
	}

	return result
}

// 円
type CircleCollision struct {
	Sprite *Sprite
	Center Vector2
	Radius float64 // 半径
}

func NewCircleCollision(sprite *Sprite, pos Vector2, radius float64) *CircleCollision {
	return &CircleCollision{Sprite: sprite, Center: pos, Radius: radius}
}

func (c *CircleCollision) Test(o CollisionTester) bool {
	result := false
	switch v := o.(type) {
	case *PolygonCollision:
		result = TestCirclePolygon(c, v)
	case *CircleCollision:
		result = TestCircleCircle(c, v)
	case *CompositeCollision:
		result = TestCircleComposit(c, v)
	default:
		result = false
	}

	return result
}

// 複合形状
type CompositeCollision struct {
	Collisions []CollisionTester
	Operator   CompositeOperator // 0:or、1:and
}

func NewCompositeCollision(collisions []CollisionTester, operator CompositeOperator) *CompositeCollision {
	return &CompositeCollision{Collisions: collisions, Operator: operator}
}

func (c *CompositeCollision) Test(o CollisionTester) bool {
	result := false
	switch v := o.(type) {
	case *PolygonCollision:
		result = TestPolygonComposit(v, c)
	case *CircleCollision:
		result = TestCircleComposit(v, c)
	case *CompositeCollision:
		result = TestCompositComposit(v, c)
	default:
		result = false
	}

	return result
}

// 複合形状のAnd/Or条件
type CompositeOperator int

const (
	CompositeOr  CompositeOperator = 0
	CompositeAnd CompositeOperator = 1
)

// 点と円の判定
func TestPointCircle(pos Vector2, c1 *CircleCollision) bool {
	cpos := c1.Center.Sub(c1.Sprite.GetOrigin()).Rotate(c1.Sprite.Angle).Add(c1.Sprite.GetOrigin()).Add(c1.Sprite.GetCollisionPos()).Sub(c1.Sprite.GetOffset())
	return pos.DistanceSquaredTo(cpos) < c1.Radius*c1.Radius
}

// 点と凸型多角形の判定
func TestPointPolygon(pos Vector2, p1 *PolygonCollision) bool {
	// 外積を計算して境界内に点があるかを判定する
	r := make([]Vector2, 0, len(p1.Vertices)+1)
	for _, p := range p1.Vertices {
		// 各頂点から回転原点を引いてから回転、回転原点とベース座標を足すことでグローバル座標を算出する
		v := p.Sub(p1.Sprite.GetOrigin()).Rotate(p1.Sprite.Angle).Add(p1.Sprite.GetOrigin()).Add(p1.Sprite.GetCollisionPos()).Sub(p1.Sprite.GetOffset())
		r = append(r, v)
	}

	// 1個目の頂点をスライスに追加する
	// これでr[n+1]-r[n]がエッジのベクトルになる
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
func TestCircleCircle(c1 *CircleCollision, c2 *CircleCollision) bool {
	cpos1 := c1.Center.Sub(c1.Sprite.GetOrigin()).Rotate(c1.Sprite.Angle).Add(c1.Sprite.GetOrigin()).Add(c1.Sprite.GetCollisionPos()).Sub(c1.Sprite.GetOffset())
	cpos2 := c2.Center.Sub(c2.Sprite.GetOrigin()).Rotate(c2.Sprite.Angle).Add(c2.Sprite.GetOrigin()).Add(c2.Sprite.GetCollisionPos()).Sub(c2.Sprite.GetOffset())
	return cpos1.DistanceSquaredTo(cpos2) < (c1.Radius+c2.Radius)*(c1.Radius+c2.Radius)
}

// 円と凸型多角形の判定
func TestCirclePolygon(c1 *CircleCollision, p1 *PolygonCollision) bool {
	cpos := c1.Center.Sub(c1.Sprite.GetOrigin()).Rotate(c1.Sprite.Angle).Add(c1.Sprite.GetOrigin()).Add(c1.Sprite.GetCollisionPos()).Sub(c1.Sprite.GetOffset())

	r := make([]Vector2, 0, len(p1.Vertices)+1)
	for _, v := range p1.Vertices {
		// 各頂点から回転原点を引いてから回転、回転原点とベース座標を足すことでグローバル座標を算出する
		v := v.Sub(p1.Sprite.GetOrigin()).Rotate(p1.Sprite.Angle).Add(p1.Sprite.GetOrigin()).Add(p1.Sprite.GetCollisionPos()).Sub(p1.Sprite.GetOffset())
		r = append(r, v)
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
func TestPolygonPolygon(c1 *PolygonCollision, c2 *PolygonCollision) bool {
	r1 := make([]Vector2, 0, len(c1.Vertices)+1)
	for _, p := range c1.Vertices {
		// 各頂点から回転原点を引いてから回転、回転原点とベース座標を足すことでグローバル座標を算出する
		v := p.Sub(c1.Sprite.GetOrigin()).Rotate(c1.Sprite.Angle).Add(c1.Sprite.GetOrigin()).Add(c1.Sprite.GetCollisionPos()).Sub(c1.Sprite.GetOffset())
		r1 = append(r1, v)
	}

	// 1個目の頂点をスライスに追加する
	// これでr1[n+1]-r1[n]がエッジのベクトルになる
	r1 = append(r1, r1[0])

	r2 := make([]Vector2, 0, len(c2.Vertices)+1)
	for _, p := range c2.Vertices {
		// 各頂点から回転原点を引いてから回転、回転原点とベース座標を足すことでグローバル座標を算出する
		v := p.Sub(c2.Sprite.GetOrigin()).Rotate(c2.Sprite.Angle).Add(c2.Sprite.GetOrigin()).Add(c2.Sprite.GetCollisionPos()).Sub(c2.Sprite.GetOffset())
		r2 = append(r2, v)
	}

	// 1個目の頂点をスライスに追加する
	// これでr2[n+1]-r2[n]がエッジのベクトルになる
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
func TestPointComposit(pos Vector2, co *CompositeCollision) bool {
	result := false
	for _, d := range co.Collisions {
		switch v := d.(type) {
		case *PolygonCollision:
			result = TestPointPolygon(pos, v)
		case *CircleCollision:
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
func TestCircleComposit(c *CircleCollision, co *CompositeCollision) bool {
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

func TestPolygonComposit(p *PolygonCollision, co *CompositeCollision) bool {
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

func TestCompositComposit(c1 *CompositeCollision, c2 *CompositeCollision) bool {
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
