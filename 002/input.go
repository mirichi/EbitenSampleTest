package main

// ブラウザでスマホを想定した入力ロジック
// マウスとタッチの両方に対応する
// ・マウスのボタンは左のみ
// ・タッチは最初の1つのみを対象とする
// ・座標はマウスとタッチで共通、押したときだけ更新する

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		panic(err)
	}
	mplusFaceSource = s
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

func drawtext(screen *ebiten.Image, msg string, y int) {
	const (
		normalFontSize = 24
		bigFontSize    = 48
	)

	op := &text.DrawOptions{}
	op.GeoM.Translate(20, float64(y))

	text.Draw(screen, msg, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op)
}

func Input_Draw(screen *ebiten.Image) {
	// Draw info
	drawtext(screen, fmt.Sprintf(" TouchID:%d\n", touchID), 20)
	drawtext(screen, fmt.Sprintf(" TouchJustPressed:%v\n", touchJustPressed), 40)
	drawtext(screen, fmt.Sprintf(" TouchPressed:%v\n", touchPressed), 60)
	drawtext(screen, fmt.Sprintf(" TouchJustReleased:%v\n", touchJustReleased), 80)
	drawtext(screen, fmt.Sprintf(" MouseJustPressed:%v\n", mouseJustPressed), 100)
	drawtext(screen, fmt.Sprintf(" MousePressed:%v\n", mousePressed), 120)
	drawtext(screen, fmt.Sprintf(" MouseJustReleased:%v\n", mouseJustReleased), 140)
	drawtext(screen, fmt.Sprintf(" X:%d\n", x), 160)
	drawtext(screen, fmt.Sprintf(" Y:%d\n", y), 180)
	drawtext(screen, fmt.Sprintf(" OldX:%d\n", oldX), 200)
	drawtext(screen, fmt.Sprintf(" OldY:%d\n", oldY), 220)

}
