package objects

import (
	"MyProject/gameparts"
	"MyProject/parts"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite は画像を表示する機能を持つEntity
// EntityBaseを埋め込んで座標や親子関係の機能を再利用する
type Sprite struct {
	parts.EntityBase
	gameparts.Transform
	gameparts.Drawable
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
	s.InitEntityBase(s, x, y)
	s.InitTransform(s)
	s.InitDrawable(s, img)

	s.OriginX = float64(img.Bounds().Dx()) / 2
	s.OriginY = float64(img.Bounds().Dy()) / 2
}

// GetGlobalPos は、グローバル行列を原点(0,0)に適用し、グローバルなX, Y座標を返します。
func (s *Sprite) GetGlobalPos() (float64, float64) {
	m := s.GetGlobalMatrix()
	return m.Apply(0, 0)
}
