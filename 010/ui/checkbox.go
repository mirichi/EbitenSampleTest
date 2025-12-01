package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Checkboxはチェック状態を持つコントロール
type Checkbox struct {
	*parts.ControlBase
	*parts.MouseInteraction
	*parts.Drawable
	*parts.TextDrawable
	*parts.Focusable

	Checked        bool
	OnCheckChanged func(bool)
}

// Checkbox生成
func NewCheckbox(x, y, w, h int, text string, size int, initialChecked bool) *Checkbox {
	c := &Checkbox{
		Checked: initialChecked,
	}

	c.ControlBase = parts.NewControlBase(c, x, y, w, h)

	// クリック時の動作
	c.MouseInteraction = parts.NewMouseInteraction(c)
	c.OnClick = func() {
		c.Checked = !c.Checked
		if c.OnCheckChanged != nil {
			c.OnCheckChanged(c.Checked)
		}
	}

	// 背景とチェックマークを描画するカスタム描画関数
	c.Drawable = parts.NewDrawable(c, c.drawCheckbox)
	c.TextDrawable = parts.NewTextDrawable(c, text, size, parts.AlignLeft, parts.AlignCenter, h+5, 0, color.White, true)
	c.Focusable = parts.NewFocusable(c)

	return c
}

func (c *Checkbox) drawCheckbox(screen *ebiten.Image) {
	gx, gy := c.GetGlobalPos()
	x := float32(gx)
	y := float32(gy)
	h := float32(c.Height)

	// ボックスの枠
	vector.StrokeRect(screen, x, y, h, h, 2, color.White, false)

	// チェックされている場合の中身
	if c.Checked {
		vector.FillRect(screen, x+4, y+4, h-8, h-8, color.RGBA{0x00, 0xff, 0x00, 0xff}, false)
	}
}
