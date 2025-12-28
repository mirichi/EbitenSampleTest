package objects

import (
	"MyProject/parts"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite は画像を表示する機能を持つEntity
// EntityBaseを埋め込んで座標や親子関係の機能を再利用する
type Sprite struct {
	parts.EntityBase

	Angle float64

	// 各パラメータをX,Yに分離
	OriginX float64
	OriginY float64
	OffsetX float64
	OffsetY float64

	Image *ebiten.Image
}

// NewSprite は新しいSpriteを作成する
func NewSprite(x, y float64, img *ebiten.Image) *Sprite {
	s := &Sprite{}
	s.InitSprite(x, y, img)
	return s
}

// InitSprite はSpriteを初期化する
func (s *Sprite) InitSprite(x, y float64, img *ebiten.Image) {
	// EntityBaseを初期化
	// 第一引数は自分自身(Entity)を渡す
	s.EntityBase.InitEntityBase(s, x, y)

	s.Image = img

	// 描画関数をメインのDrawフェーズに登録
	s.AddDrawFunction(s.draw)
}

// 登録用の描画関数
func (s *Sprite) draw(screen *ebiten.Image) {
	if s.Image == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// 親の座標変更も反映された絶対座標を取得
	gx, gy := s.GetGlobalPos()

	// 座標を設定
	op.GeoM.Translate(-s.OriginX, -s.OriginY)
	op.GeoM.Rotate(s.Angle)
	op.GeoM.Translate(s.OriginX, s.OriginY)
	op.GeoM.Translate(gx, gy)
	op.GeoM.Translate(-s.OffsetX, -s.OffsetY)

	screen.DrawImage(s.Image, op)
}

// GetCollisionPos は衝突判定用のグローバル座標を返す
func (s *Sprite) GetCollisionPos() Vector2 {
	x, y := s.GetGlobalPos()
	return Vector2{X: x, Y: y}
}

// GetOrigin は回転中心座標をVector2で返す(衝突判定用)
func (s *Sprite) GetOrigin() Vector2 {
	return Vector2{X: s.OriginX, Y: s.OriginY}
}

// GetOffset はオフセット座標をVector2で返す(衝突判定用)
func (s *Sprite) GetOffset() Vector2 {
	return Vector2{X: s.OffsetX, Y: s.OffsetY}
}
