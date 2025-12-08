package input

// ブラウザでスマホを想定した入力ロジック
// マウスとタッチの両方に対応する
// マウスのボタンは左のみ

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var x, y int = 0, 0
var oldX, oldY int = 0, 0

// タッチ情報はここで保持する
var touches []Touch = []Touch{}
var mouseTouch MouseTouch = MouseTouch{id: -1}

type Touch interface {
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

type TouchTouch struct {
	id       ebiten.TouchID
	released bool
}

func (t *TouchTouch) Pos() (int, int) {
	return ebiten.TouchPosition(t.id)
}
func (t *TouchTouch) OldPos() (int, int) {
	if t.IsJustPressed() {
		return t.Pos()
	}
	return inpututil.TouchPositionInPreviousTick(t.id)
}
func (t *TouchTouch) IsJustPressed() bool {
	return inpututil.TouchPressDuration(t.id) == 1
}
func (t *TouchTouch) IsPressed() bool {
	return inpututil.TouchPressDuration(t.id) > 0
}
func (t *TouchTouch) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.id)
}
func (t *TouchTouch) ID() ebiten.TouchID {
	return t.id
}
func (t *TouchTouch) isReleased() bool {
	return t.released
}
func (t *TouchTouch) release() {
	t.released = true
}
func (t *TouchTouch) clear() {
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

func Update() {
	oldX, oldY = x, y
	x, y = ebiten.CursorPosition()

	// 終了したタッチをスライスから削除
	n := []Touch{}
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
		touches = append(touches, &TouchTouch{id: tid, released: false})
	}

	// 新規クリックをスライスに追加
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseTouch.released = false
		touches = append(touches, &mouseTouch)
	}
}

// タッチ中の情報を返す
func GetAllTouches() []Touch {
	return touches
}

// 1個目のタッチの情報を返す
func GetFirstTouch() Touch {
	if len(touches) > 0 {
		return touches[0]
	}
	return nil
}

// マウス情報取得
func GetMouseTouch() Touch {
	return &mouseTouch
}

// ポインタ情報取得（タッチがあればそれを、なければマウス座標を返す）
// クリックしていなくてもマウス座標を返すため、ホバー判定に使用できる
func GetPointer() Touch {
	if len(touches) > 0 {
		return touches[0]
	}
	return &mouseTouch
}

// 右クリック関連
func IsRightJustPressed() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
}

func IsRightPressed() bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
}

func IsRightJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight)
}
