package main

// EventType はステージイベントの種類を表す
type EventType int

const (
	EventTypeSpawnEnemy EventType = iota
	EventTypeWait
	EventTypeSpawnBoss
	EventTypeEnd
)

// StageEvent はステージ上で発生する1つのイベントを表す
type StageEvent struct {
	Type      EventType
	Time      int       // イベント発生までの待ち時間（フレーム数）。Waitの場合に使用。
	EnemyType EnemyType // SpawnEnemyの場合の敵の種類
	X         float64   // SpawnEnemyの場合のX座標
	Y         float64   // SpawnEnemyの場合のY座標
}

// StageData はステージ全体のデータを表す
type StageData struct {
	Events []StageEvent
}

// CreateStage1 はステージ1のデータを生成して返す
func CreateStage1() *StageData {
	events := []StageEvent{}

	// Helper to add spawn event
	spawn := func(t EnemyType, x float64) {
		events = append(events, StageEvent{
			Type:      EventTypeSpawnEnemy,
			EnemyType: t,
			X:         x,
			Y:         -50, // 画面外から
		})
	}

	// Helper to add wait event
	wait := func(frames int) {
		events = append(events, StageEvent{
			Type: EventTypeWait,
			Time: frames,
		})
	}

	// Start
	wait(60)

	// Wave 1: Straight enemies
	spawn(EnemyTypeStraight, 100)
	wait(30)
	spawn(EnemyTypeStraight, 300)
	wait(30)
	spawn(EnemyTypeStraight, 500)
	wait(60)

	// Wave 2: Wave enemies
	spawn(EnemyTypeWave, 200)
	spawn(EnemyTypeWave, 400)
	wait(120)

	// Bonus: PowerUp Opportunity
	spawn(EnemyTypeMedium, 320)
	wait(180)

	// Wave 3: Fast enemies
	spawn(EnemyTypeFast, 150)
	wait(20)
	spawn(EnemyTypeFast, 320)
	wait(20)
	spawn(EnemyTypeFast, 490)
	wait(180)

	// Wave 4: Mixed
	spawn(EnemyTypeStraight, 100)
	spawn(EnemyTypeWave, 500)
	wait(60)
	spawn(EnemyTypeFast, 320)
	wait(120)

	// Boss Spawn
	events = append(events, StageEvent{
		Type: EventTypeSpawnBoss,
	})

	// End
	events = append(events, StageEvent{
		Type: EventTypeEnd,
	})

	return &StageData{Events: events}
}

// CreateStage2 はステージ2のデータを生成して返す
// ステージ1より難易度高め
func CreateStage2() *StageData {
	events := []StageEvent{}

	// Helper to add spawn event
	spawn := func(t EnemyType, x float64) {
		events = append(events, StageEvent{
			Type:      EventTypeSpawnEnemy,
			EnemyType: t,
			X:         x,
			Y:         -50, // 画面外から
		})
	}

	// Helper to add wait event
	wait := func(frames int) {
		events = append(events, StageEvent{
			Type: EventTypeWait,
			Time: frames,
		})
	}

	// Start
	wait(60)

	// Wave 1: Simultaneous Fast and Warning Shot
	spawn(EnemyTypeFast, 100)
	spawn(EnemyTypeFast, 540)
	wait(40)
	spawn(EnemyTypeFast, 200)
	spawn(EnemyTypeFast, 440)
	wait(40)
	spawn(EnemyTypeFast, 320)
	wait(60)

	// Wave 2: Wave Rush
	spawn(EnemyTypeWave, 50)
	wait(20)
	spawn(EnemyTypeWave, 150)
	wait(20)
	spawn(EnemyTypeWave, 250)
	wait(20)
	spawn(EnemyTypeWave, 350)
	wait(20)
	spawn(EnemyTypeWave, 450)
	wait(20)
	spawn(EnemyTypeWave, 550)
	wait(120)

	// Bonus: PowerUp Opportunity
	spawn(EnemyTypeMedium, 320)
	wait(180)

	// Wave 3: Mixed Hell
	spawn(EnemyTypeStraight, 50)
	spawn(EnemyTypeStraight, 590)
	spawn(EnemyTypeFast, 320)
	wait(60)
	spawn(EnemyTypeWave, 100)
	spawn(EnemyTypeWave, 540)
	spawn(EnemyTypeStraight, 320)
	wait(120)

	// Boss Spawn
	events = append(events, StageEvent{
		Type: EventTypeSpawnBoss,
	})

	// End
	events = append(events, StageEvent{
		Type: EventTypeEnd,
	})

	return &StageData{Events: events}
}

// CreateStage3 はステージ3のデータを生成して返す
// 弾幕地獄の入り口
func CreateStage3() *StageData {
	events := []StageEvent{}

	spawn := func(t EnemyType, x float64) {
		events = append(events, StageEvent{
			Type:      EventTypeSpawnEnemy,
			EnemyType: t,
			X:         x,
			Y:         -50,
		})
	}

	wait := func(frames int) {
		events = append(events, StageEvent{
			Type: EventTypeWait,
			Time: frames,
		})
	}

	wait(60)

	// Wave 1: Introduction to Shooting Enemies
	spawn(EnemyTypeShooting, 200)
	wait(60)
	spawn(EnemyTypeShooting, 440)
	wait(60)
	spawn(EnemyTypeShooting, 320)
	wait(120)

	// Wave 2: Shooting Rush
	spawn(EnemyTypeShooting, 100)
	spawn(EnemyTypeShooting, 540)
	wait(60)
	spawn(EnemyTypeShooting, 220)
	spawn(EnemyTypeShooting, 420)
	wait(60)
	spawn(EnemyTypeMedium, 320) // PowerUp
	wait(180)

	// Wave 3: Mixed Assault
	spawn(EnemyTypeFast, 50)
	spawn(EnemyTypeFast, 590)
	spawn(EnemyTypeShooting, 320)
	wait(60)
	spawn(EnemyTypeWave, 100)
	spawn(EnemyTypeWave, 540)
	spawn(EnemyTypeShooting, 200)
	spawn(EnemyTypeShooting, 440)
	wait(120)

	// Boss Spawn
	events = append(events, StageEvent{
		Type: EventTypeSpawnBoss,
	})

	events = append(events, StageEvent{
		Type: EventTypeEnd,
	})

	return &StageData{Events: events}
}
