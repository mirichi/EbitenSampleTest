package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"MyProject/input"
	"MyProject/objects"
)

type GameScene struct {
	Root              *objects.Container
	Player            *Player
	EnemyGroup        *objects.Container
	BulletGroup       *objects.Container
	EnemyBulletGroup  *objects.Container
	ItemGroup         *objects.Container // Added
	GameOver          bool
	Score             int
	KillCount         int
	BossSpawned       bool
	Boss              Enemy // Changed from *Boss
	ShootTimer        int
	EffectManager     *EffectManager
	BossDefeatedTimer int
	StageStartTimer   int // Added

	// Stage Management
	Stage             *StageData
	StageNum          int
	CurrentEventIndex int
	WaitTimer         int
	StageFinished     bool
}

func NewGameScene() *GameScene {
	gs := &GameScene{}
	gs.InitGameScene()
	return gs
}

func (gs *GameScene) InitGameScene() {
	// Root Container
	gs.Root = objects.NewContainer(0, 0)

	// Groups for Enemies and Bullets
	gs.EnemyGroup = objects.NewContainer(0, 0)
	gs.Root.AddChild(gs.EnemyGroup)

	gs.BulletGroup = objects.NewContainer(0, 0)
	gs.Root.AddChild(gs.BulletGroup)

	// Group for Enemy Bullets
	gs.EnemyBulletGroup = objects.NewContainer(0, 0)
	gs.Root.AddChild(gs.EnemyBulletGroup)

	// Group for Items
	gs.ItemGroup = objects.NewContainer(0, 0)
	gs.Root.AddChild(gs.ItemGroup)

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

	// Stage Init
	gs.StageNum = 1
	gs.LoadStage(gs.StageNum)
}

func (gs *GameScene) LoadStage(stageNum int) {
	gs.StageNum = stageNum

	switch stageNum {
	case 1:
		gs.Stage = CreateStage1()
	case 2:
		gs.Stage = CreateStage2()
	default:
		// Fallback or Loop
		gs.Stage = CreateStage1()
		gs.StageNum = 1
	}

	gs.CurrentEventIndex = 0
	gs.WaitTimer = 0
	gs.StageFinished = false
	gs.BossSpawned = false
	gs.Boss = nil
	gs.BossDefeatedTimer = 0
	gs.StageStartTimer = 180 // Show title for 3 seconds

	// Clear existing enemies for safety (if any remain)
	// Usually handled by cleanup or boss dead logic, but good to be sure if loading mid-game
}

func (gs *GameScene) Update() error {
	// Reset logic
	if gs.GameOver {
		if input.IsActionJustPressed(input.ActionShot) ||
			input.IsActionJustPressed(input.ActionShowMenu) ||
			input.IsActionJustPressed(input.ActionReset) {
			gs.InitGameScene()
		}
		return nil
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
			// Stage Clear! Next Stage!
			gs.LoadStage(gs.StageNum + 1)
			return nil
		}
	}

	if gs.StageStartTimer > 0 {
		gs.StageStartTimer--
	}

	// Transfer Pending Bullets from Boss
	if gs.BossSpawned && gs.Boss != nil {
		bullets := gs.Boss.GetPendingBullets()
		if len(bullets) > 0 {
			for _, b := range bullets {
				gs.EnemyBulletGroup.AddChild(b)
			}
		}
	}

	gs.Root.Update()

	// Shooting Logic (Auto-Fire)
	if gs.ShootTimer > 0 {
		gs.ShootTimer--
	}
	isShooting := input.IsActionPressed(input.ActionShot)
	if t := input.GetPointer(); t != nil && t.IsPressed() {
		isShooting = true
	}

	if isShooting {
		if gs.ShootTimer <= 0 {
			switch gs.Player.PowerLevel {
			case 1:
				bullet := NewBullet(gs.Player.X(), gs.Player.Y()-24)
				gs.BulletGroup.AddChild(bullet)
			case 2:
				b1 := NewBullet(gs.Player.X()-8, gs.Player.Y()-24)
				b2 := NewBullet(gs.Player.X()+8, gs.Player.Y()-24)
				gs.BulletGroup.AddChild(b1)
				gs.BulletGroup.AddChild(b2)
			case 3:
				b1 := NewBulletAngled(gs.Player.X()-16, gs.Player.Y()-24, -math.Pi/2-0.2, 8.0)
				b2 := NewBullet(gs.Player.X(), gs.Player.Y()-24)
				b3 := NewBulletAngled(gs.Player.X()+16, gs.Player.Y()-24, -math.Pi/2+0.2, 8.0)
				gs.BulletGroup.AddChild(b1)
				gs.BulletGroup.AddChild(b2)
				gs.BulletGroup.AddChild(b3)
			case 4:
				b1 := NewBullet(gs.Player.X()-8, gs.Player.Y()-24)
				b2 := NewBullet(gs.Player.X()+8, gs.Player.Y()-24)
				b3 := NewBulletAngled(gs.Player.X()-20, gs.Player.Y()-24, -math.Pi/2-0.2, 8.0)
				b4 := NewBulletAngled(gs.Player.X()+20, gs.Player.Y()-24, -math.Pi/2+0.2, 8.0)
				gs.BulletGroup.AddChild(b1)
				gs.BulletGroup.AddChild(b2)
				gs.BulletGroup.AddChild(b3)
				gs.BulletGroup.AddChild(b4)
			default: // Fallback same as 4
				b1 := NewBullet(gs.Player.X()-8, gs.Player.Y()-24)
				b2 := NewBullet(gs.Player.X()+8, gs.Player.Y()-24)
				b3 := NewBulletAngled(gs.Player.X()-20, gs.Player.Y()-24, -math.Pi/2-0.2, 8.0)
				b4 := NewBulletAngled(gs.Player.X()+20, gs.Player.Y()-24, -math.Pi/2+0.2, 8.0)
				gs.BulletGroup.AddChild(b1)
				gs.BulletGroup.AddChild(b2)
				gs.BulletGroup.AddChild(b3)
				gs.BulletGroup.AddChild(b4)
			}
			gs.ShootTimer = 8 // Cooldown frames
		}
	}

	// Stage Progression
	if !gs.StageFinished && !gs.BossSpawned {
		if gs.WaitTimer > 0 {
			gs.WaitTimer--
		} else {
			// Process Events
			for gs.CurrentEventIndex < len(gs.Stage.Events) {
				event := gs.Stage.Events[gs.CurrentEventIndex]
				gs.CurrentEventIndex++

				switch event.Type {
				case EventTypeSpawnEnemy:
					enemy := SpawnEnemyByType(event.EnemyType, event.X, event.Y)
					gs.EnemyGroup.AddChild(enemy)
					// Continue to next event immediately (unless wait is explicitly called)
				case EventTypeWait:
					gs.WaitTimer = event.Time
					goto BreakLoop // Wait commands pause processing
				case EventTypeSpawnBoss:
					gs.BossSpawned = true
					if gs.StageNum == 2 {
						gs.Boss = NewBoss2(320, -100)
					} else {
						gs.Boss = NewBoss(320, -100) // Start off-screen
					}
					gs.EnemyGroup.AddChild(gs.Boss)
					goto BreakLoop
				case EventTypeEnd:
					gs.StageFinished = true
					goto BreakLoop
				}
			}
		BreakLoop:
		}
	}

	// Collision Logic via Children
	gs.checkCollisions()

	// Cleanup dead entities
	gs.cleanup()

	return nil
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
						enemy.ApplyDamage(1)
						if enemy.IsDead() {
							gs.KillCount++
							gs.Score += 100

							isBoss := false
							if gs.Boss != nil && enemy == gs.Boss {
								isBoss = true
							}

							if !isBoss {
								ex, ey := enemy.GetGlobalPos()
								gs.EffectManager.SpawnExplosion(ex, ey, color.White)

								// Drop Item
								if item := enemy.DropItem(); item != nil {
									gs.ItemGroup.AddChild(item)
								}
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

	// Items vs Player
	iChildren := gs.ItemGroup.Children
	for _, iChild := range iChildren {
		if item, ok := iChild.(*Item); ok {
			if !item.IsDead() {
				if gs.Player.Test(item) {
					item.MarkDead()
					gs.Player.PowerUp()
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

	// Remove dead items
	iChildren := gs.ItemGroup.Children
	for i := len(iChildren) - 1; i >= 0; i-- {
		child := iChildren[i]
		if item, ok := child.(*Item); ok {
			if item.IsDead() {
				gs.ItemGroup.RemoveChild(item)
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

	// Stage Title
	if gs.StageStartTimer > 0 {
		// Blink effect or just solid? Solid for now.
		msg := fmt.Sprintf("STAGE %d", gs.StageNum)
		ebitenutil.DebugPrintAt(screen, msg, 290, 240)
	}

	// HUD
	// Use GUI instead
	// ebitenutil.DebugPrint(screen, "Arrow Keys: Move, Z: Shoot")
}
