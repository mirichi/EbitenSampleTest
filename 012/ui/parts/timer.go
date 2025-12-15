package parts

// Timerは指定時間経過後に処理を実行するパーツ
// 画面には表示されないが、partsとして追加することでUpdateが呼ばれて時間が進む
type Timer struct {
	Control Control

	Interval int  // タイムアウトまでの時間（フレーム数）
	Current  int  // 現在の経過時間
	Repeat   bool // 繰り返し実行するか
	Running  bool // タイマーが動作中か

	OnTimer func() // タイムアウト時に呼ばれる関数
}

// Timer生成
func NewTimer(c Control, interval int, repeat bool) *Timer {
	t := &Timer{}
	t.InitTimer(c, interval, repeat)
	return t
}

func (t *Timer) InitTimer(c Control, interval int, repeat bool) {
	t.Control = c
	t.Interval = interval
	t.Repeat = repeat

	// Update処理の登録
	c.GetControlBase().AddUpdateFunction(t.updateTimer)
}

func (t *Timer) updateTimer() {
	if !t.Running {
		return
	}

	t.Current++
	if t.Current >= t.Interval {
		if t.OnTimer != nil {
			t.OnTimer()
		}

		if t.Repeat {
			t.Current = 0
		} else {
			t.Running = false
			t.Current = 0
		}
	}
}

// Start はタイマーを開始する
func (t *Timer) Start() {
	t.Running = true
	t.Current = 0
}

// Stop はタイマーを停止する
func (t *Timer) Stop() {
	t.Running = false
	t.Current = 0
}

// Resume はタイマーを再開する
func (t *Timer) Resume() {
	t.Running = true
}

// Suspendはタイマーを中断する
func (t *Timer) Suspend() {
	t.Running = false
}
