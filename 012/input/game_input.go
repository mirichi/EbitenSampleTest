package input

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Action はゲーム内の抽象的な操作を表すID
type Action int

const (
	ActionNone Action = iota
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionShowMenu
	ActionShot
	ActionFast
)

// キーボード割り当て
var KeyBindings = map[Action][]ebiten.Key{
	ActionMoveUp:    {ebiten.KeyUp, ebiten.KeyW},
	ActionMoveDown:  {ebiten.KeyDown, ebiten.KeyS},
	ActionMoveLeft:  {ebiten.KeyLeft, ebiten.KeyA},
	ActionMoveRight: {ebiten.KeyRight, ebiten.KeyD},
	ActionShot:      {ebiten.KeyZ, ebiten.KeySpace},
	ActionFast:      {ebiten.KeyShift},
	ActionShowMenu:  {ebiten.KeyEscape},
}

// ゲームパッドボタン割り当て（汎用的な配置を想定）
var GamepadButtonBindings = map[Action][]ebiten.GamepadButton{
	ActionMoveUp:    {ebiten.GamepadButton10},                       // D-pad Up
	ActionMoveDown:  {ebiten.GamepadButton12},                       // D-pad Down
	ActionMoveLeft:  {ebiten.GamepadButton13},                       // D-pad Left
	ActionMoveRight: {ebiten.GamepadButton11},                       // D-pad Right
	ActionShot:      {ebiten.GamepadButton0, ebiten.GamepadButton1}, // A or B (bottom/right face buttons)
	ActionFast:      {ebiten.GamepadButton4, ebiten.GamepadButton5}, // L1/R1
	ActionShowMenu:  {ebiten.GamepadButton6, ebiten.GamepadButton7}, // Start/Select
}

// ゲームパッドスティック設定
// スティック入力をデジタルな方向入力として扱うための閾値
const stickThreshold = 0.5

var currentGamepadID ebiten.GamepadID = -1

// UpdateGamepadID は接続されているゲームパッドを確認し、IDを更新する
// input.Update() から呼び出すことを想定
func UpdateGamepadID() {
	// 接続されているパッドがなければ再検索
	// 現在のIDが無効になっているかどうかもチェックすべきだが、
	// 簡易的に毎回（あるいは定期的に）接続リストを確認して先頭のものを採用する戦略をとる
	ids := ebiten.AppendGamepadIDs(nil)
	if len(ids) > 0 {
		currentGamepadID = ids[0]
	} else {
		currentGamepadID = -1
	}
}

// IsActionPressed は指定されたアクションが入力されているかを返す
func IsActionPressed(action Action) bool {
	// キーボード判定
	if keys, ok := KeyBindings[action]; ok {
		for _, k := range keys {
			if ebiten.IsKeyPressed(k) {
				return true
			}
		}
	}

	// ゲームパッド判定
	if currentGamepadID >= 0 {
		// ボタン判定
		if btns, ok := GamepadButtonBindings[action]; ok {
			for _, b := range btns {
				if ebiten.IsGamepadButtonPressed(currentGamepadID, b) {
					return true
				}
			}
		}

		// スティック判定（移動系アクションのみ）
		// 左スティック (Axis 0, 1) を想定
		switch action {
		case ActionMoveLeft:
			if ebiten.GamepadAxisValue(currentGamepadID, 0) < -stickThreshold {
				return true
			}
		case ActionMoveRight:
			if ebiten.GamepadAxisValue(currentGamepadID, 0) > stickThreshold {
				return true
			}
		case ActionMoveUp:
			if ebiten.GamepadAxisValue(currentGamepadID, 1) < -stickThreshold {
				return true
			}
		case ActionMoveDown:
			if ebiten.GamepadAxisValue(currentGamepadID, 1) > stickThreshold {
				return true
			}
		}
	}

	return false
}

// IsActionJustPressed は指定されたアクションが押された瞬間かを返す
func IsActionJustPressed(action Action) bool {
	// キーボード判定
	if keys, ok := KeyBindings[action]; ok {
		for _, k := range keys {
			if inpututil.IsKeyJustPressed(k) {
				return true
			}
		}
	}

	// ゲームパッド判定
	if currentGamepadID >= 0 {
		// ボタン判定
		if btns, ok := GamepadButtonBindings[action]; ok {
			for _, b := range btns {
				if inpututil.IsGamepadButtonJustPressed(currentGamepadID, b) {
					return true
				}
			}
		}

		// スティックのJustPressed判定は少し複雑（前回値の保持が必要）なので
		// 今回はボタンのみ対応とする。移動にJustPressedを使うことは稀なため。
	}

	return false
}

// GetGamepadDebugInfo はデバッグ用に現在のゲームパッドの状態を文字列で返します
func GetGamepadDebugInfo() string {
	if currentGamepadID < 0 {
		return "No Gamepad"
	}

	msg := fmt.Sprintf("ID: %d, Name: %s\nButtons: ", currentGamepadID, ebiten.GamepadName(currentGamepadID))

	// ボタン (標準的なボタン数30くらいまで確認)
	for b := ebiten.GamepadButton(0); b <= ebiten.GamepadButtonMax; b++ {
		if ebiten.IsGamepadButtonPressed(currentGamepadID, b) {
			msg += fmt.Sprintf("%d ", b)
		}
	}
	msg += "\nAxes: "

	// 軸
	numAxes := ebiten.GamepadAxisCount(currentGamepadID)
	for a := 0; a < numAxes; a++ {
		v := ebiten.GamepadAxisValue(currentGamepadID, a)
		if math.Abs(v) > 0.1 { // ノイズ除去
			msg += fmt.Sprintf("%d:%.2f ", a, v)
		}
	}

	return msg
}
