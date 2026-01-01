package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"MyProject/input"
	"MyProject/objects"
	"MyProject/parts"
)

type GameScene struct {
	parts.EntityBase
	Root              *objects.Container
	Player            *Player
	EnemyGroup        *objects.Container
	BulletGroup       *objects.Container
	EnemyBulletGroup  *objects.Container
	GameOver          bool
	Score             int
	KillCount         int
	BossSpawned       bool
	Boss              *Boss
	ShootTimer        int
	EffectManager     *EffectManager
	BossDefeatedTimer int
}

func NewGameScene() *GameScene {
	gs := &GameScene{}
	gs.InitGameScene()
	return gs
}

func (gs *GameScene) InitGameScene() {
	// Initialize EntityBase
	gs.EntityBase.InitEntityBase(gs, 0, 0)

	// Create Root Container
	gs.Root = objects.NewContainer(0, 0)

	// Groups for Enemies and Bullets
	gs.EnemyGroup = objects.NewContainer(0, 0)
	gs.Root.AddChild(gs.EnemyGroup)

	gs.BulletGroup = objects.NewContainer(0, 0)
	gs.Root.AddChild(gs.BulletGroup)

	// Group for Enemy Bullets
	gs.EnemyBulletGroup = objects.NewContainer(0, 0)
	gs.Root.AddChild(gs.EnemyBulletGroup)

	// Player
	gs.Player = NewPlayer()
	gs.Root.AddChild(gs.Player)

	gs.GameOver = false
	gs.Score = 0
	gs.KillCount = 0
	gs.BossSpawned = false
	gs.Boss = nil
	gs.ShootTimer = 0
	gs.EffectManager = NewEffectManager()
	gs.Root.AddChild(gs.EffectManager)
	gs.BossDefeatedTimer = 0
}

func (gs *GameScene) Update() {
	// Reset logic
	if gs.GameOver {
		if input.IsActionJustPressed(input.ActionShot) || input.IsActionJustPressed(input.ActionShowMenu) {
			gs.InitGameScene()
		}
		return
	}
	// Reset logic (Boss Defeated)
	if gs.BossSpawned && gs.Boss != nil && gs.Boss.IsDead() {
		// Start timer if not started
		if gs.BossDefeatedTimer == 0 {
			gs.BossDefeatedTimer = 180 // 3 seconds

			// Initial Big Explosion
			ex, ey := gs.Boss.GetGlobalPos()
			gs.EffectManager.SpawnBossExplosion(ex, ey)
		}
	}

	if gs.BossDefeatedTimer > 0 {
		gs.BossDefeatedTimer--
		// Random small explosions during wait
		if gs.BossDefeatedTimer%10 == 0 && gs.BossDefeatedTimer > 60 {
			// Random position around boss
			if gs.Boss != nil {
				ex, ey := gs.Boss.GetGlobalPos()
				rx := ex - 50 + rand.Float64()*100
				ry := ey - 50 + rand.Float64()*100
				gs.EffectManager.SpawnExplosion(rx, ry, color.RGBA{255, 100, 0, 255})
			}
		}

		if gs.BossDefeatedTimer <= 0 {
			// Time to reset
			gs.InitGameScene()
			return
		}

		// Continue updating visual effects but maybe stop player or game logic?
		// For now, let everything run but just wait for reset.
	}

	// Transfer Pending Bullets from Boss
	if gs.BossSpawned && gs.Boss != nil {
		if len(gs.Boss.PendingBullets) > 0 {
			for _, b := range gs.Boss.PendingBullets {
				gs.EnemyBulletGroup.AddChild(b)
			}
			gs.Boss.PendingBullets = []*EnemyBullet{} // Clear
		}
	}

	gs.Root.Update()

	// Shooting Logic (Auto-Fire)
	if gs.ShootTimer > 0 {
		gs.ShootTimer--
	}
	if input.IsActionPressed(input.ActionShot) {
		if gs.ShootTimer <= 0 {
			bullet := NewBullet(gs.Player.X(), gs.Player.Y()-24)
			gs.BulletGroup.AddChild(bullet)
			gs.ShootTimer = 8 // Cooldown frames
		}
	}

	// Spawning Enemies
	if !gs.BossSpawned {
		// Spawn normal enemies until 10 kills
		if gs.KillCount < 10 {
			if rand.Float64() < 0.02 {
				enemy := SpawnEnemy()
				gs.EnemyGroup.AddChild(enemy)
			}
		} else {
			// Stop spawning, wait for clear
			if len(gs.EnemyGroup.Children) == 0 {
				// Spawn Boss
				gs.BossSpawned = true
				gs.Boss = NewBoss(320, -100) // Start off-screen
				gs.EnemyGroup.AddChild(gs.Boss)
			}
		}
	}

	// Collision Logic via Children
	gs.checkCollisions()

	// Cleanup dead entities
	gs.cleanup()
}

func (gs *GameScene) checkCollisions() {
	// Enemy vs Player
	children := gs.EnemyGroup.Children
	for _, child := range children {
		if enemy, ok := child.(Enemy); ok {
			// Skip dead enemies (like Boss dying frame)
			if enemy.IsDead() {
				continue
			}

			if gs.Player.Test(enemy) {
				gs.GameOver = true
			}

			// Enemy vs Bullets
			bChildren := gs.BulletGroup.Children
			for _, bChild := range bChildren {
				if bullet, ok := bChild.(*Bullet); ok {
					if !bullet.IsDead && bullet.Test(enemy) {
						bullet.IsDead = true
						enemy.MarkDead()
						if enemy.IsDead() {
							gs.KillCount++
							gs.Score += 100
							if enemy != gs.Boss {
								ex, ey := enemy.GetGlobalPos()
								gs.EffectManager.SpawnExplosion(ex, ey, color.White)
							}
						}
						break
					}
				}
			}
		}
	}

	// Enemy Bullets vs Player
	ebChildren := gs.EnemyBulletGroup.Children
	for _, ebChild := range ebChildren {
		if bullet, ok := ebChild.(*EnemyBullet); ok {
			if !bullet.IsDead {
				if gs.Player.Test(bullet) {
					gs.GameOver = true
				}
			}
		}
	}
}

func (gs *GameScene) cleanup() {
	// Remove dead bullets
	// Iterate backwards to safely remove
	bChildren := gs.BulletGroup.Children
	for i := len(bChildren) - 1; i >= 0; i-- {
		child := bChildren[i]
		if bullet, ok := child.(*Bullet); ok {
			if bullet.IsDead {
				gs.BulletGroup.RemoveChild(bullet)
			}
		}
	}

	// Remove dead enemies
	eChildren := gs.EnemyGroup.Children
	for i := len(eChildren) - 1; i >= 0; i-- {
		child := eChildren[i]
		if enemy, ok := child.(Enemy); ok {
			if enemy.IsDead() {
				gs.EnemyGroup.RemoveChild(enemy)
			}
		}
	}

	// Remove dead enemy bullets
	ebChildren := gs.EnemyBulletGroup.Children
	for i := len(ebChildren) - 1; i >= 0; i-- {
		child := ebChildren[i]
		if bullet, ok := child.(*EnemyBullet); ok {
			if bullet.IsDead {
				gs.EnemyBulletGroup.RemoveChild(bullet)
			}
		}
	}
}

func (gs *GameScene) Draw(screen *ebiten.Image) {
	if gs.GameOver {
		// ebitenutil.DebugPrint(screen, "GAME OVER (Press R to Reset)")
		return
	}

	// Draw everything via Root
	gs.Root.Draw(screen)

	// HUD
	// Use GUI instead
	// ebitenutil.DebugPrint(screen, "Arrow Keys: Move, Z: Shoot")
}
