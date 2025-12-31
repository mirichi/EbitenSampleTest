package gameparts

import (
	"MyProject/parts"

	"github.com/hajimehoshi/ebiten/v2"
)

// GlobalMatrixer は、グローバルな変換行列を提供できるコンポーネントのためのインターフェースです。
// これを実装することで、階層的な行列計算が可能になります。
type GlobalMatrixer interface {
	GetGlobalMatrix() ebiten.GeoM
	GetVertexMatrix() ebiten.GeoM
}

// Transform は、拡大縮小(ScaleX, ScaleY)、回転(Angle)などの座標系機能を提供します。
// 位置(X, Y)はEntityから直接取得します。
// Entity内にコンポーネントとして埋め込むか、Entityをラップして使用することを想定しています。
type Transform struct {
	Entity parts.Entity // このTransformが属する、またはラップしているEntity
	// X, Y is managed by Entity
	ScaleX  float64
	ScaleY  float64
	Angle   float64 // 回転角度（ラジアン）
	OriginX float64
	OriginY float64
}

// NewTransform は、スケールを1で初期化した新しいTransformコンポーネントを作成します。
func NewTransform(e parts.Entity) *Transform {
	eb := &Transform{}
	eb.InitTransform(e)
	return eb
}

func (eb *Transform) InitTransform(e parts.Entity) {
	eb.Entity = e
	eb.ScaleX = 1
	eb.ScaleY = 1
	eb.Angle = 0
	eb.OriginX = 0
	eb.OriginY = 0
}

// GetGlobalMatrix は、最終的なグローバル変換行列を計算します。
// ローカル変換（拡大縮小 -> 回転 -> 移動）と、親のグローバル変換を合成します。
// 親がGlobalMatrixerを実装している場合は行列を連結し、そうでない場合は単純な座標移動を行います。
func (t *Transform) GetGlobalMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}

	// ローカル変換の適用順序: 拡大縮小 -> 回転 -> 移動
	m.Scale(t.ScaleX, t.ScaleY)
	m.Rotate(t.Angle)

	// 位置はEntityから取得
	eb := t.Entity.GetEntityBase()
	m.Translate(eb.X, eb.Y)

	// 親の変換を適用
	if eb.Parent != nil {
		parent := eb.Parent

		// 親がグローバル行列をサポートしているか確認
		if gm, ok := parent.(GlobalMatrixer); ok {
			pm := gm.GetGlobalMatrix()
			m.Concat(pm)
		} else {
			// フォールバック: 親がグローバル座標（単純な移動）しか持っていない場合
			px, py := parent.GetGlobalPos()
			m.Translate(px, py)
		}
	}

	return m
}

// GetVertexMatrix は、頂点座標用の変換行列を計算します。
func (t *Transform) GetVertexMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-t.OriginX, -t.OriginY)
	m.Concat(t.GetGlobalMatrix())
	return m
}
