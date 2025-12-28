package ui

import (
	"MyProject/parts"
	"image/color"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Checkboxはチェック状態を持つControl
type Checkbox struct {
	InteractiveControl
	parts.TextDrawable
	parts.Focusable

	Checked           bool
	OnCheckChanged    func(bool)
	animationProgress float64
}

// Checkbox生成
func NewCheckbox(x, y, h float64, text string, size float64, initialChecked bool) *Checkbox {
	c := &Checkbox{}
	c.InitCheckbox(x, y, h, text, size, initialChecked)
	return c
}

func (c *Checkbox) InitCheckbox(x, y, h float64, text string, size float64, initialChecked bool) {
	w := float64(int(float64(h) * 1.8))
	c.InitInteractiveControl(x, y, w, h)
	c.InitTextDrawable(c, text, size, parts.AlignLeft, parts.AlignCenter, w+5, 0, color.White, true)
	c.InitFocusable(c)
	c.OnDraw = c.drawCheckbox
	c.Checked = initialChecked

	if c.Checked {
		c.animationProgress = 1.0
	} else {
		c.animationProgress = 0.0
	}

	c.AddAfterUpdateFunction(func() {
		target := 0.0
		if c.Checked {
			target = 1.0
		}
		// シンプルな線形補間でアニメーション (0.2はスピード調整)
		c.animationProgress += (target - c.animationProgress) * 0.2
		if math.Abs(target-c.animationProgress) < 0.01 {
			c.animationProgress = target
		}
	})

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
	w := float32(c.Width)
	h := float32(c.Height)

	theme := parts.CurrentTheme

	// 背景色をアニメーションに合わせて補間
	offColor := theme.CheckOffColor
	onColor := theme.CheckOnColor

	// シンプルに進行度0.5を境に切り替える
	bgColor := offColor
	if c.animationProgress > 0.5 {
		bgColor = onColor
	}

	// カプセル背景の描画
	r := h / 2
	// 左円
	vector.FillCircle(screen, x+r, y+r, r, bgColor, true)
	// 右円
	vector.FillCircle(screen, x+w-r, y+r, r, bgColor, true)
	// 真ん中Rect
	vector.FillRect(screen, x+r, y, w-h, h, bgColor, false)

	// ノブ（白い丸）
	knobRadius := (h - 4) / 2
	knobXStart := x + 2 + knobRadius
	knobXEnd := x + w - 2 - knobRadius

	currentKnobX := knobXStart + (knobXEnd-knobXStart)*float32(c.animationProgress)

	vector.FillCircle(screen, currentKnobX, y+r, knobRadius, color.White, true)
}
