package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Checkboxはチェック状態を持つコントロール
type Checkbox struct {
	InteractiveWidget
	parts.TextDrawable
	parts.Focusable

	Checked        bool
	OnCheckChanged func(bool)
}

// Checkbox生成
func NewCheckbox(x, y, w, h int, text string, size int, initialChecked bool) *Checkbox {
	c := &Checkbox{}
	c.InitCheckbox(nil, x, y, w, h, text, size, initialChecked)
	return c
}

func (c *Checkbox) InitCheckbox(g parts.AddChilder, x, y, w, h int, text string, size int, initialChecked bool) {
	c.InitInteractiveWidget(nil, x, y, w, h)
	c.InitTextDrawable(c, text, size, parts.AlignLeft, parts.AlignCenter, h+5, 0, color.White, true)
	c.InitFocusable(c)
	c.OnDraw = c.drawCheckbox
	c.Checked = initialChecked
	if g != nil {
		g.AddChild(c)
	}

	// クリック時の動作
	c.OnClick = func() {
		c.Checked = !c.Checked
		if c.OnCheckChanged != nil {
			c.OnCheckChanged(c.Checked)
		}
	}
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
