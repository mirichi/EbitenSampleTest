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

	OriginX float64
	OriginY float64

	ColorScale ebiten.ColorScale

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
	s.InitEntityBase(s, x, y)

	s.Image = img
	s.OriginX = float64(img.Bounds().Dx()) / 2
	s.OriginY = float64(img.Bounds().Dy()) / 2

	s.ColorScale.Scale(1, 1, 1, 1)

	// 描画関数をメインのDrawフェーズに登録
	s.AddDrawFunction(s.draw)
}

func (s *Sprite) GetGlobalMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}

	// Local Transform: Translate(-Origin) -> Rotate -> Translate(X, Y)
	// (X, Y) becomes the Pivot Point in parent space
	m.Translate(-s.OriginX, -s.OriginY)
	m.Rotate(s.Angle)
	m.Translate(s.X, s.Y)

	// Parent Transform
	if p := s.EntityBase.Parent; p != nil {
		if gm, ok := p.(GlobalMatrixer); ok {
			pm := gm.GetGlobalMatrix()
			m.Concat(pm)
		} else {
			px, py := p.GetGlobalPos()
			m.Translate(px, py)
		}
	}
	return m
}

// 登録用の描画関数
func (s *Sprite) draw(screen *ebiten.Image) {
	if s.Image == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// get Global Matrix
	m := s.GetGlobalMatrix()

	op.GeoM = m
	op.ColorScale = s.ColorScale

	screen.DrawImage(s.Image, op)
}

// GetGlobalPos override
func (s *Sprite) GetGlobalPos() (float64, float64) {
	m := ebiten.GeoM{}
	m.Translate(s.X, s.Y)

	// Parent Transform
	if p := s.EntityBase.Parent; p != nil {
		if gm, ok := p.(GlobalMatrixer); ok {
			pm := gm.GetGlobalMatrix()
			m.Concat(pm)
		} else {
			px, py := p.GetGlobalPos()
			m.Translate(px, py)
		}
	}
	return m.Apply(0, 0)
}
