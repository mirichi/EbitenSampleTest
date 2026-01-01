package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/objects"
)

type EnemyType int

const (
	EnemyTypeStraight EnemyType = iota
	EnemyTypeWave
	EnemyTypeFast
	EnemyTypeMedium
)

// SpawnEnemyByType は指定されたタイプと位置で敵を生成する
func SpawnEnemyByType(t EnemyType, x, y float64) Enemy {
	var img *ebiten.Image
	var speed float64
	var col color.Color
	w, h := 48, 48

	switch t {
	case EnemyTypeStraight:
		col = color.RGBA{255, 0, 0, 255}
		speed = 2.0
	case EnemyTypeWave:
		col = color.RGBA{0, 255, 0, 255}
		speed = 1.5
	case EnemyTypeFast:
		col = color.RGBA{255, 255, 0, 255}
		speed = 4.0
	case EnemyTypeMedium:
		col = color.RGBA{150, 0, 255, 255} // Purple
		speed = 0.5                        // Slow
	}

	img = ebiten.NewImage(w, h)
	img.Fill(col)

	// x, y は中心座標として扱いたいが、Spriteは左上が原点とは限らない（Originによる）
	// 今回のSpriteは中心がOriginになっているので、x, yをそのまま渡せば中心になる
	s := objects.NewSprite(x, y, img)
	// 回転は一旦なしで
	// s.Angle = rand.Float64() * 6.28

	hp := 1
	if t == EnemyTypeMedium {
		hp = 20 // Hard
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
	default:
		return &StraightEnemy{EnemyBase: base}
	}
}
