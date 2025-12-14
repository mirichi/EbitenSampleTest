package parts

import "image/color"

// Theme はアプリケーション全体の色テーマを定義する
type Theme struct {
	// テキスト
	Text         color.Color
	DisabledText color.Color

	// ボタン
	ButtonNormal color.Color
	ButtonHover  color.Color
	ButtonPress  color.Color

	// ウィンドウ
	TitleBar   color.Color
	ClientArea color.Color

	// 入力コントロール
	TextBoxBackground color.Color
	TextBoxBorder     color.Color

	// チェックボックス
	CheckOffColor color.Color
	CheckOnColor  color.Color

	// スクロールバー
	ScrollBackground color.Color
	ScrollKnob       color.Color
	ScrollKnobHover  color.Color

	// ポップアップ
	PopupBackground color.Color
	PopupHover      color.Color

	// フォーカス
	FocusBorder color.Color
}

// CurrentTheme はアプリケーション全体で使用されるテーマ
var CurrentTheme = DefaultTheme()

// DefaultTheme はダークテーマのデフォルト値を返す
func DefaultTheme() *Theme {
	return &Theme{
		// テキスト
		Text:         color.White,
		DisabledText: color.RGBA{0x80, 0x80, 0x80, 0xff},

		// ボタン
		ButtonNormal: color.RGBA{0x60, 0x60, 0x60, 0xff},
		ButtonHover:  color.RGBA{0x80, 0x80, 0x80, 0xff},
		ButtonPress:  color.RGBA{0x40, 0x40, 0x40, 0xff},

		// ウィンドウ
		TitleBar:   color.RGBA{0x00, 0x40, 0x00, 0xff},
		ClientArea: color.RGBA{0x30, 0x30, 0x30, 0xff},

		// 入力コントロール
		TextBoxBackground: color.RGBA{0x20, 0x20, 0x20, 0xff},
		TextBoxBorder:     color.RGBA{0x80, 0x80, 0x80, 0xff},

		// チェックボックス
		CheckOffColor: color.RGBA{0x60, 0x60, 0x60, 0xff},
		CheckOnColor:  color.RGBA{0x34, 0xC7, 0x59, 0xff},

		// スクロールバー
		ScrollBackground: color.RGBA{0x20, 0x20, 0x20, 0xff},
		ScrollKnob:       color.RGBA{0x80, 0x80, 0x80, 0xff},
		ScrollKnobHover:  color.RGBA{0xA0, 0xA0, 0xA0, 0xff},

		// ポップアップ
		PopupBackground: color.RGBA{0x40, 0x40, 0x40, 0xff},
		PopupHover:      color.RGBA{0x60, 0x60, 0x80, 0xff},

		// フォーカス
		FocusBorder: color.RGBA{0xD0, 0xD0, 0xD0, 0xff},
	}
}
