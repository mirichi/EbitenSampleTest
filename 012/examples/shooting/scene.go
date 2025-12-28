package main

import (
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"MyProject/objects"
	"MyProject/parts"
)

type GameScene struct {
	parts.EntityBase
	Root        *objects.Container
	Player      *Player
	EnemyGroup  *objects.Container
	BulletGroup *objects.Container
	GameOver    bool
	Score       int
	KillCount   int
	BossSpawned bool
	Boss        *Boss
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

	// Player
	gs.Player = NewPlayer()
	gs.Root.AddChild(gs.Player)

	gs.GameOver = false
	gs.Score = 0
	gs.KillCount = 0
	gs.BossSpawned = false
	gs.Boss = nil
}

func (gs *GameScene) Update() {
	if gs.GameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			gs.InitGameScene()
		}
		return
	}

	// Reset logic (Boss Defeated)
	if gs.BossSpawned && gs.Boss != nil && gs.Boss.IsDead() {
		// Simple reset for now, or you could show "YOU WIN"
		// Let's just reset the scene for loop
		gs.InitGameScene()
		return
	}

	gs.Root.Update()

	// Explicitly handle Shooting input here
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		bullet := NewBullet(gs.Player.X()+12, gs.Player.Y()-20)
		gs.BulletGroup.AddChild(bullet)
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
			if len(gs.EnemyGroup.GetChildren()) == 0 {
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
	children := gs.EnemyGroup.GetChildren()
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
			bChildren := gs.BulletGroup.GetChildren()
			for _, bChild := range bChildren {
				if bullet, ok := bChild.(*Bullet); ok {
					if !bullet.IsDead && bullet.Test(enemy) {
						bullet.IsDead = true
						enemy.MarkDead()
						if enemy.IsDead() {
							gs.KillCount++
							gs.Score += 100
						}
						break
					}
				}
			}
		}
	}
}

func (gs *GameScene) cleanup() {
	// Remove dead bullets
	// Iterate backwards to safely remove
	bChildren := gs.BulletGroup.GetChildren()
	for i := len(bChildren) - 1; i >= 0; i-- {
		child := bChildren[i]
		if bullet, ok := child.(*Bullet); ok {
			if bullet.IsDead {
				gs.BulletGroup.RemoveChild(bullet)
			}
		}
	}

	// Remove dead enemies
	eChildren := gs.EnemyGroup.GetChildren()
	for i := len(eChildren) - 1; i >= 0; i-- {
		child := eChildren[i]
		if enemy, ok := child.(Enemy); ok {
			if enemy.IsDead() {
				gs.EnemyGroup.RemoveChild(enemy)
			}
		}
	}
}

func (gs *GameScene) Draw(screen *ebiten.Image) {
	if gs.GameOver {
		ebitenutil.DebugPrint(screen, "GAME OVER (Press R to Reset)")
		return
	}

	// Draw everything via Root
	gs.Root.Draw(screen)

	// HUD
	ebitenutil.DebugPrint(screen, "Arrow Keys: Move, Z: Shoot")
	ebitenutil.DebugPrintAt(screen, "Score: "+strconv.Itoa(gs.Score), 0, 15)
	if gs.BossSpawned && gs.Boss != nil {
		ebitenutil.DebugPrintAt(screen, "BOSS HP: "+strconv.Itoa(gs.Boss.HP), 0, 30)
	}
}
