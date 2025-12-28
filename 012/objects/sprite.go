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

	ScaleX float64
	ScaleY float64

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
	s.ScaleX = 1
	s.ScaleY = 1

	// No strict identity needed for ColorScale as zero value allows Set, but explicit init is clearer
	// s.ColorScale is zero value, which is effectively identity-like for ScaleWithColor but standard Scale is 0?
	// Ebiten 2.5+: ColorScale zero value is identity? No.
	// We should probably rely on not setting it if it's default, or init to 1,1,1,1
	// Wait, Ebiten's ColorScale usage:
	// op.ColorScale.Scale(r,g,b,a).
	// If we want it to be 1,1,1,1 initially.
	// Actually, let's just leave it zero value, and we apply it.
	// If it is zero value, it might be 0,0,0,0?
	// Let's check `ebiten.ColorScale`.
	// It's a struct (r,g,b,a float32). Zero value is all 0.
	// So we must initialize it to identity (1,1,1,1).
	s.ColorScale.Scale(1, 1, 1, 1)

	// 描画関数をメインのDrawフェーズに登録
	s.AddDrawFunction(s.draw)
}

// GlobalMatrixer is defined in container.go, but we can rely on interface matching or duplicate locally if needed.
// To avoid strict dependency on container.go (though they are in same package), we just define it locally or use it if exported.
// Note: In Go, they are in the same package 'objects', so 'GlobalMatrixer' is visible if it was exported.
// But I didn't export it in Container (I named it 'GlobalMatrixer' with uppercase G). good.

func (s *Sprite) GetGlobalMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}

	// Local Transform: Translate(-Origin) -> Scale -> Rotate -> Translate(Origin) -> Translate(X, Y)
	m.Translate(-s.OriginX, -s.OriginY)
	m.Scale(s.ScaleX, s.ScaleY)
	m.Rotate(s.Angle)
	m.Translate(s.OriginX, s.OriginY)
	m.Translate(-s.OffsetX, -s.OffsetY)
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
	m := s.GetGlobalMatrix()
	return m.Apply(0, 0)
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
