package parts

import "image/color"

// Theme はアプリケーション全体の色テーマを定義する
type Theme struct {
	// テキスト
	Text color.Color

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
	CheckBoxBackground color.Color
	CheckMark          color.Color

	// スクロールバー
	ScrollBackground color.Color
	ScrollKnob       color.Color

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
		Text: color.White,

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
		CheckBoxBackground: color.RGBA{0x40, 0x40, 0x40, 0xff},
		CheckMark:          color.White,

		// スクロールバー
		ScrollBackground: color.RGBA{0x20, 0x20, 0x20, 0xff},
		ScrollKnob:       color.RGBA{0x60, 0x60, 0x60, 0xff},

		// ポップアップ
		PopupBackground: color.RGBA{0x40, 0x40, 0x40, 0xff},
		PopupHover:      color.RGBA{0x60, 0x60, 0x80, 0xff},

		// フォーカス
		FocusBorder: color.RGBA{0xD0, 0xD0, 0xD0, 0xff},
	}
}

// LightTheme はライトテーマを返す(AIが作ったので実用的じゃないっぽい)
func LightTheme() *Theme {
	return &Theme{
		// テキスト
		Text: color.Black,

		// ボタン
		ButtonNormal: color.RGBA{0xE0, 0xE0, 0xE0, 0xff},
		ButtonHover:  color.RGBA{0xC8, 0xC8, 0xC8, 0xff},
		ButtonPress:  color.RGBA{0xB0, 0xB0, 0xB0, 0xff},

		// ウィンドウ
		TitleBar:   color.RGBA{0x00, 0x80, 0x00, 0xff},
		ClientArea: color.RGBA{0xF0, 0xF0, 0xF0, 0xff},

		// 入力コントロール
		TextBoxBackground: color.White,
		TextBoxBorder:     color.RGBA{0x80, 0x80, 0x80, 0xff},

		// チェックボックス
		CheckBoxBackground: color.RGBA{0xE0, 0xE0, 0xE0, 0xff},
		CheckMark:          color.Black,

		// スクロールバー
		ScrollBackground: color.RGBA{0xE0, 0xE0, 0xE0, 0xff},
		ScrollKnob:       color.RGBA{0xA0, 0xA0, 0xA0, 0xff},

		// ポップアップ
		PopupBackground: color.RGBA{0xF0, 0xF0, 0xF0, 0xff},
		PopupHover:      color.RGBA{0xC0, 0xC0, 0xE0, 0xff},

		// フォーカス
		FocusBorder: color.RGBA{0x40, 0x40, 0x40, 0xff},
	}
}
