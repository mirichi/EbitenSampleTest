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

// 標準ゲームパッドボタン割り当て
var StandardGamepadButtonBindings = map[Action][]ebiten.StandardGamepadButton{
	ActionMoveUp:    {ebiten.StandardGamepadButtonLeftTop},
	ActionMoveDown:  {ebiten.StandardGamepadButtonLeftBottom},
	ActionMoveLeft:  {ebiten.StandardGamepadButtonLeftLeft},
	ActionMoveRight: {ebiten.StandardGamepadButtonLeftRight},
	ActionShot:      {ebiten.StandardGamepadButtonRightBottom, ebiten.StandardGamepadButtonRightRight},     // A/B, Cross/Circle
	ActionFast:      {ebiten.StandardGamepadButtonFrontTopLeft, ebiten.StandardGamepadButtonFrontTopRight}, // L1/R1, LB/RB
	ActionShowMenu:  {ebiten.StandardGamepadButtonCenterRight},                                             // Start/Menu
}

// ゲームパッドスティック設定
// スティック入力をデジタルな方向入力として扱うための閾値
const stickThreshold = 0.5

var currentGamepadID ebiten.GamepadID = -1

// UpdateGamepadID は接続されているゲームパッドを確認し、IDを更新する
// input.Update() から呼び出すことを想定
func UpdateGamepadID() {
	// 接続されているパッドがなければ再検索
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
		// Standard Gamepad Button Checks
		if btns, ok := StandardGamepadButtonBindings[action]; ok {
			for _, b := range btns {
				if ebiten.IsStandardGamepadButtonPressed(currentGamepadID, b) {
					return true
				}
			}
		}

		// Standard Gamepad Stick Checks (Left Stick)
		switch action {
		case ActionMoveLeft:
			if ebiten.StandardGamepadAxisValue(currentGamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal) < -stickThreshold {
				return true
			}
		case ActionMoveRight:
			if ebiten.StandardGamepadAxisValue(currentGamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal) > stickThreshold {
				return true
			}
		case ActionMoveUp:
			if ebiten.StandardGamepadAxisValue(currentGamepadID, ebiten.StandardGamepadAxisLeftStickVertical) < -stickThreshold {
				return true
			}
		case ActionMoveDown:
			if ebiten.StandardGamepadAxisValue(currentGamepadID, ebiten.StandardGamepadAxisLeftStickVertical) > stickThreshold {
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
		// Standard Gamepad Button Checks
		if btns, ok := StandardGamepadButtonBindings[action]; ok {
			for _, b := range btns {
				if inpututil.IsStandardGamepadButtonJustPressed(currentGamepadID, b) {
					return true
				}
			}
		}
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

	// Standard Gamepad Info
	msg += "\nStandard Buttons: "
	for b := ebiten.StandardGamepadButton(0); b <= ebiten.StandardGamepadButtonMax; b++ {
		if ebiten.IsStandardGamepadButtonPressed(currentGamepadID, b) {
			msg += fmt.Sprintf("%d ", b)
		}
	}

	return msg
}
