package main

// ブラウザでスマホを想定した入力ロジック
// マウスとタッチの両方に対応する
// ・マウスのボタンは左のみ
// ・タッチは最初の1つのみを対象とする
// ・座標はマウスとタッチで共通、押したときだけ更新する

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var touchID ebiten.TouchID = -1
var touchJustPressed bool = false
var touchPressed bool = false
var touchJustReleased bool = false
var mouseJustPressed bool = false
var mousePressed bool = false
var mouseJustReleased bool = false
var x, y int = 0, 0
var oldX, oldY int = 0, 0

// タッチ情報のポインタはここで保持する
var touches map[ebiten.TouchID]*Touch = map[ebiten.TouchID]*Touch{}
var mouseTouch MouseTouch = MouseTouch{id: -1}

type TouchInfo interface {
	Pos() (int, int)
	OldPos() (int, int)
	IsJustPressed() bool
	IsPressed() bool
	IsJustReleased() bool
	ID() ebiten.TouchID
}

type Touch struct {
	id ebiten.TouchID
}

func (t *Touch) Pos() (int, int) {
	return ebiten.TouchPosition(t.id)
}
func (t *Touch) OldPos() (int, int) {
	return inpututil.TouchPositionInPreviousTick(t.id)
}
func (t *Touch) IsJustPressed() bool {
	return inpututil.TouchPressDuration(t.id) == 1
}
func (t *Touch) IsPressed() bool {
	return inpututil.TouchPressDuration(t.id) > 0
}
func (t *Touch) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.id)
}
func (t *Touch) ID() ebiten.TouchID {
	return t.id
}

type MouseTouch struct {
	id ebiten.TouchID
}

func (t *MouseTouch) Pos() (int, int) {
	return ebiten.CursorPosition()
}
func (t *MouseTouch) OldPos() (int, int) {
	return oldX, oldY
}
func (t *MouseTouch) IsJustPressed() bool {
	return mouseJustPressed
}
func (t *MouseTouch) IsPressed() bool {
	return mousePressed
}
func (t *MouseTouch) IsJustReleased() bool {
	return mouseJustReleased
}
func (t *MouseTouch) ID() ebiten.TouchID {
	return t.id
}

func init() {
}

func Input_Update() {
	touchJustPressed = false
	touchJustReleased = false

	// 前回タッチしていた
	if touchID != -1 {
		if inpututil.TouchPressDuration(touchID) == 0 {
			// タッチが終了した
			touchID = -1
			touchPressed = false
			touchJustReleased = true
		} else {
			// タッチ継続中
			oldX, oldY = x, y
			x, y = ebiten.TouchPosition(touchID)
		}
	}

	// タッチが開始されたらtouchIDをセットする
	// 追加のタッチは無視
	if touchID == -1 {
		touchIDs := inpututil.AppendJustPressedTouchIDs(nil)
		if len(touchIDs) > 0 {
			touchID = touchIDs[0]
			touchJustPressed = true
			touchPressed = true
			oldX, oldY = x, y
			x, y = ebiten.TouchPosition(touchID)
		}
	}

	mouseJustPressed = false
	mouseJustReleased = false

	// マウスの左ボタンのみ対応
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if mousePressed {
			mouseJustReleased = true
			mousePressed = false
		}
	} else {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mouseJustPressed = true
			mousePressed = true
		}
		oldX, oldY = x, y
		x, y = ebiten.CursorPosition()
	}

	// 終了したタッチをマップから削除
	for tid, touch := range touches {
		// Todo:離された次のフレームに別のタッチでIDが使いまわされるとバグる気がする
		if inpututil.TouchPressDuration(tid) == 0 && !inpututil.IsTouchJustReleased(tid) {
			touch.id = -1
			delete(touches, tid)
		}
	}
	// 新規タッチをマップに登録
	for _, tid := range inpututil.AppendJustPressedTouchIDs(nil) {
		touches[tid] = &Touch{id: tid}
	}
}

// touchesとmouseTouchを統合したスライスを返す
func AllTouches() []TouchInfo {
	result := make([]TouchInfo, 0, len(touches)+1)
	for _, t := range touches {
		result = append(result, t)
	}
	result = append(result, &mouseTouch)
	return result
}

func IsButtonJustPressed() bool {
	return touchJustPressed || mouseJustPressed
}

func IsButtonPressed() bool {
	return touchPressed || mousePressed
}

func IsButtonJustReleased() bool {
	return touchJustReleased || mouseJustReleased
}

func CurrentPos() (int, int) {
	return x, y
}

func OldPos() (int, int) {
	return oldX, oldY
}
