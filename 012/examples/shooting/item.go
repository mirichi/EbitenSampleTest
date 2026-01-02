package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
)

type ItemType int

const (
	ItemTypePowerUp ItemType = iota
)

type Item struct {
	*objects.Sprite
	*objects.CircleCollider
	Type   ItemType
	isDead bool
}

func NewItem(x, y float64, t ItemType) *Item {
	var resID ResourceID
	switch t {
	case ItemTypePowerUp:
		resID = ImgItemPowerUp
	}
	img := GetResourceManager().GetImage(resID)

	s := objects.NewSprite(x, y, img)

	i := &Item{
		Sprite: s,
		Type:   t,
		isDead: false,
	}
	// コリジョンは少し大きめに
	i.CircleCollider = objects.NewCircleCollider(s, objects.Vector2{X: 12, Y: 12}, 20)

	return i
}

func (i *Item) Update() {
	i.Sprite.Y += 2.0 // ゆっくり落ちる

	if i.Sprite.Y > ScreenHeight+20 {
		i.isDead = true
	}

	i.Sprite.Update()
}

func (i *Item) Draw(screen *ebiten.Image) {
	if !i.isDead {
		i.Sprite.Draw(screen)
	}
}

func (i *Item) IsDead() bool {
	return i.isDead
}

func (i *Item) MarkDead() {
	i.isDead = true
}

// ItemはEntityとして振る舞える (Spriteを埋め込んでいるため)
