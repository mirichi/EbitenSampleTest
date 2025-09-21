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
var justPressed bool = false
var pressed bool = false
var justReleased bool = false
var x, y int = 0, 0
var oldX, oldY int = 0, 0

func init() {
}

func Input_Update() {
	justPressed = false
	justReleased = false

	// 前回タッチしていた
	if touchID != -1 {
		if inpututil.TouchPressDuration(touchID) == 0 {
			// タッチが終了した
			touchID = -1
			pressed = false
			justReleased = true
		} else {
			// タッチ継続中
			oldX, oldY = x, y
			x, y = ebiten.TouchPosition(touchID)
		}
	}

	// タッチが開始されたらtouchIDをセットする
	// 追加のタッチは無視
	if touchID == -1 {
		touchIDs := inpututil.AppendJustReleasedTouchIDs(nil)
		if len(touchIDs) > 0 {
			touchID = touchIDs[0]
			justPressed = true
			pressed = true
		}
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if pressed {
			justReleased = true
			pressed = false
		}
	} else {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			justPressed = true
			pressed = true
		}
		oldX, oldY = x, y
		x, y = ebiten.CursorPosition()
	}
}

func IsButtonJustPressed() bool {
	return justPressed
}

func IsButtonPressed() bool {
	return pressed
}

func IsButtonJustReleased() bool {
	return justReleased
}

func CurrectPos() (int, int) {
	return x, y
}

func OldPos() (int, int) {
	return oldX, oldY
}
