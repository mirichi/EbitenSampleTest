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

func CurrectPos() (int, int) {
	return x, y
}

func OldPos() (int, int) {
	return oldX, oldY
}
