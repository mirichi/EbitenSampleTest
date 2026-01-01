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
