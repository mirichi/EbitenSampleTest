package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// ResourceID is a type for resource identifiers.
type ResourceID string

const (
	// Enemy Images
	ImgEnemyStraight ResourceID = "img_enemy_straight"
	ImgEnemyWave     ResourceID = "img_enemy_wave"
	ImgEnemyFast     ResourceID = "img_enemy_fast"
	ImgEnemyMedium   ResourceID = "img_enemy_medium"
	ImgEnemyShooting ResourceID = "img_enemy_shooting"

	// Boss Images
	ImgBossBody    ResourceID = "img_boss_body"
	ImgBossHand    ResourceID = "img_boss_hand"
	ImgBoss2Body   ResourceID = "img_boss2_body"
	ImgBoss2Turret ResourceID = "img_boss2_turret"
	ImgBoss3Core   ResourceID = "img_boss3_core"
	ImgBoss3Wing   ResourceID = "img_boss3_wing"
	ImgBoss3Cannon ResourceID = "img_boss3_cannon"

	// Bullet Images
	ImgBulletPlayer ResourceID = "img_bullet_player"
	ImgBulletEnemy  ResourceID = "img_bullet_enemy"
	ImgBulletBossP  ResourceID = "img_bullet_boss_p" // Pink/Default?

	// Item Images
	ImgItemPowerUp ResourceID = "img_item_powerup"

	// Player Images
	ImgPlayer ResourceID = "img_player"
)

// ResourceManager manages game assets.
type ResourceManager struct {
	images map[ResourceID]*ebiten.Image
}

var instance *ResourceManager

// GetResourceManager returns the singleton instance of ResourceManager.
func GetResourceManager() *ResourceManager {
	if instance == nil {
		instance = &ResourceManager{
			images: make(map[ResourceID]*ebiten.Image),
		}
		instance.loadDefaultAssets()
	}
	return instance
}

// GetImage returns the image for the given ID.
func (rm *ResourceManager) GetImage(id ResourceID) *ebiten.Image {
	return rm.images[id]
}

// loadDefaultAssets generates procedural assets (rectangles/fills) for now.
func (rm *ResourceManager) loadDefaultAssets() {
	// --- Enemies ---
	rm.registerFillImage(ImgEnemyStraight, 48, 48, color.RGBA{255, 0, 0, 255})
	rm.registerFillImage(ImgEnemyWave, 48, 48, color.RGBA{0, 255, 0, 255})
	rm.registerFillImage(ImgEnemyFast, 48, 48, color.RGBA{255, 255, 0, 255})
	rm.registerFillImage(ImgEnemyMedium, 48, 48, color.RGBA{150, 0, 255, 255})
	rm.registerFillImage(ImgEnemyShooting, 48, 48, color.RGBA{255, 120, 0, 255})

	// --- Boss 1 ---
	rm.registerFillImage(ImgBossBody, 96, 96, color.RGBA{200, 50, 50, 255})
	rm.registerFillImage(ImgBossHand, 36, 36, color.RGBA{150, 150, 0, 255})

	// --- Boss 2 ---
	rm.registerFillImage(ImgBoss2Body, 120, 80, color.RGBA{50, 50, 200, 255})
	rm.registerFillImage(ImgBoss2Turret, 40, 40, color.RGBA{100, 200, 255, 255})

	// --- Boss 3 ---
	rm.registerFillImage(ImgBoss3Core, 64, 64, color.RGBA{100, 0, 150, 255})
	rm.registerFillImage(ImgBoss3Wing, 32, 96, color.RGBA{80, 0, 120, 255})
	rm.registerFillImage(ImgBoss3Cannon, 24, 48, color.RGBA{200, 50, 200, 255})

	// --- Bullets ---
	rm.registerFillImage(ImgBulletPlayer, 12, 12, color.RGBA{255, 255, 0, 255})
	rm.registerFillImage(ImgBulletEnemy, 12, 12, color.RGBA{255, 0, 100, 255})

	// --- Items ---
	rm.registerFillImage(ImgItemPowerUp, 24, 24, color.RGBA{255, 100, 100, 255})

	// --- Player ---
	rm.registerFillImage(ImgPlayer, 48, 48, color.RGBA{0, 0, 255, 255})
}

func (rm *ResourceManager) registerFillImage(id ResourceID, w, h int, clr color.Color) {
	img := ebiten.NewImage(w, h)
	img.Fill(clr)
	rm.images[id] = img
}
