package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
)

type EnemyType int

const (
	EnemyTypeStraight EnemyType = iota
	EnemyTypeWave
	EnemyTypeFast
	EnemyTypeMedium
	EnemyTypeShooting
)

// SpawnEnemyByType は指定されたタイプと位置で敵を生成する
func SpawnEnemyByType(t EnemyType, x, y float64) Enemy {
	var img *ebiten.Image
	// Determine Resource ID and parameters
	var resID ResourceID
	var speed float64

	switch t {
	case EnemyTypeStraight:
		resID = ImgEnemyStraight
		speed = 2.0
	case EnemyTypeWave:
		resID = ImgEnemyWave
		speed = 1.5
	case EnemyTypeFast:
		resID = ImgEnemyFast
		speed = 4.0
	case EnemyTypeMedium:
		resID = ImgEnemyMedium
		speed = 0.5
	case EnemyTypeShooting:
		resID = ImgEnemyShooting
		speed = 1.0
	}

	// Retrieve shared image from ResourceManager
	img = GetResourceManager().GetImage(resID)

	// x, y は中心座標として扱いたいが、Spriteは左上が原点とは限らない（Originによる）
	// 今回のSpriteは中心がOriginになっているので、x, yをそのまま渡せば中心になる
	s := objects.NewSprite(x, y, img)
	// 回転は一旦なしで
	// s.Angle = rand.Float64() * 6.28

	hp := 1
	if t == EnemyTypeMedium {
		hp = 20 // Hard
	}
	if t == EnemyTypeShooting {
		hp = 3
	}

	base := EnemyBase{
		Sprite: s,
		isDead: false,
		baseX:  x,
		speedY: speed,
		hp:     hp,
	}
	base.PolygonCollider = objects.NewRectCollider(s)

	switch t {
	case EnemyTypeStraight:
		return &StraightEnemy{EnemyBase: base}
	case EnemyTypeWave:
		return &WaveEnemy{EnemyBase: base}
	case EnemyTypeFast:
		return &FastEnemy{EnemyBase: base}
	case EnemyTypeMedium:
		return &MediumEnemy{EnemyBase: base}
	case EnemyTypeShooting:
		return &ShootingEnemy{EnemyBase: base, PendingBullets: []*EnemyBullet{}}
	default:
		return &StraightEnemy{EnemyBase: base}
	}
}
