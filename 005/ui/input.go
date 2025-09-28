package ui

// ブラウザでスマホを想定した入力ロジック
// マウスとタッチの両方に対応する
// ・マウスのボタンは左のみ
// ・タッチは最初の1つのみを対象とする
// ・座標はマウスとタッチで共通、押したときだけ更新する

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var x, y int = 0, 0
var oldX, oldY int = 0, 0

// タッチ情報はここで保持する
var touches []TouchInfo = []TouchInfo{}
var mouseTouch MouseTouch = MouseTouch{id: -1}

type TouchInfo interface {
	Pos() (int, int)
	OldPos() (int, int)
	IsJustPressed() bool
	IsPressed() bool
	IsJustReleased() bool
	ID() ebiten.TouchID
	isReleased() bool
	release()
	clear()
}

type Touch struct {
	id       ebiten.TouchID
	released bool
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
func (t *Touch) isReleased() bool {
	return t.released
}
func (t *Touch) release() {
	t.released = true
}
func (t *Touch) clear() {
	t.id = -1
}

type MouseTouch struct {
	id       ebiten.TouchID
	released bool
}

func (t *MouseTouch) Pos() (int, int) {
	return ebiten.CursorPosition()
}
func (t *MouseTouch) OldPos() (int, int) {
	return oldX, oldY
}
func (t *MouseTouch) IsJustPressed() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}
func (t *MouseTouch) IsPressed() bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}
func (t *MouseTouch) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}
func (t *MouseTouch) ID() ebiten.TouchID {
	return t.id
}
func (t *MouseTouch) isReleased() bool {
	return t.released
}
func (t *MouseTouch) release() {
	t.released = true
}
func (t *MouseTouch) clear() {
}

func init() {
}

func Input_Update() {
	oldX, oldY = x, y
	x, y = ebiten.CursorPosition()

	// 終了したタッチをスライスから削除
	n := []TouchInfo{}
	for _, t := range touches {
		if t.isReleased() {
			// 削除対象のタッチをクリアしておく
			t.clear()
		} else {
			// 削除対象じゃないタッチは残す
			n = append(n, t)

			// 離されたタッチは次回の削除対象となる
			if t.IsJustReleased() {
				t.release()
			}
		}
	}
	touches = n

	// 新規タッチをスライスに追加
	for _, tid := range inpututil.AppendJustPressedTouchIDs(nil) {
		touches = append(touches, &Touch{id: tid, released: false})
	}

	// 新規クリックをスライスに追加
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		touches = append(touches, &mouseTouch)
	}
}

// タッチ中の情報を返す
func AllTouches() []TouchInfo {
	return touches
}

// 1個目のタッチの情報を返す
func FirstTouch() TouchInfo {
	if len(touches) > 0 {
		return touches[0]
	}
	return nil
}
