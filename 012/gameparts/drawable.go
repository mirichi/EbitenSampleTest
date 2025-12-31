package gameparts

import (
	"MyProject/parts"

	"github.com/hajimehoshi/ebiten/v2"
)

// Drawable は、画像の描画機能を提供するコンポーネントです。
// Entityの位置(GlobalMatrix)に合わせて描画を行います。
type Drawable struct {
	Entity     parts.Entity
	Image      *ebiten.Image
	ColorScale ebiten.ColorScale
}

// NewDrawable は、新しいDrawableコンポーネントを作成します。
func NewDrawable(e parts.Entity, img *ebiten.Image) *Drawable {
	eb := &Drawable{}
	eb.InitDrawable(e, img)
	return eb
}

func (eb *Drawable) InitDrawable(e parts.Entity, img *ebiten.Image) {
	eb.Entity = e
	eb.Image = img
	eb.ColorScale.Scale(1, 1, 1, 1)
	eb.Entity.GetEntityBase().AddBeforeDrawFunction(eb.draw)
}

// draw は、登録用の描画関数です。Groupingなどの管理下で呼ばれることを想定しています。
func (s *Drawable) draw(screen *ebiten.Image) {
	if s.Image == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	m := ebiten.GeoM{}

	if gm, ok := s.Entity.(GlobalMatrixer); ok {
		pm := gm.GetVertexMatrix()
		m.Concat(pm)
	} else {
		// GlobalMatrixer非対応の場合は、単純な位置のみ適用
		px, py := s.Entity.GetGlobalPos()
		m.Translate(px, py)
	}

	op.GeoM = m
	op.ColorScale = s.ColorScale

	screen.DrawImage(s.Image, op)
}
