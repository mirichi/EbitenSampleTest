package objects

import (
	"MyProject/parts"

	"github.com/hajimehoshi/ebiten/v2"
)

// ContainerSprite は画像を表示する機能を持つContainer
// EntityBaseを埋め込んで座標や親子関係の機能を再利用する
// Spriteを継承することで、Spriteとしての振る舞い（衝突判定への利用など）も可能にする
type ContainerSprite struct {
	Sprite
	parts.Grouping
}

// NewContainerSprite は新しいContainerSpriteを作成する
func NewContainerSprite(x, y float64, img *ebiten.Image) *ContainerSprite {
	s := &ContainerSprite{}
	s.InitContainerSprite(x, y, img)
	return s
}

// InitContainerSprite はContainerSpriteを初期化する
func (s *ContainerSprite) InitContainerSprite(x, y float64, img *ebiten.Image) {
	// Spriteを初期化
	s.InitSprite(x, y, img)

	// Grouping初期化
	// コンテナ自身を親として登録
	s.InitGrouping(s)
}
