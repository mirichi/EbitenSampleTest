package ui

import (
	"MyProject/ui/parts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// 値をクランプする
func clampValue(value float64, allRange, viewRange int) float64 {
	if value < 0 {
		return 0
	}
	max := float64(allRange - viewRange)
	if value > max {
		return max
	}
	return value
}

// スライダーのツマミ位置を計算
func calculateKnobPos(value float64, sliderSize, knobSize, viewRange, allRange int) int {
	if allRange > viewRange {
		return int(value * float64(sliderSize-knobSize) / float64(allRange-viewRange))
	}
	return 0
}

// ScrollButtonはスクロールバー用ボタン
type ScrollButton struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.TextDrawable
}

func (b *ScrollButton) InitScrollButton(x, y, w, h int, text string, size int) {
	theme := parts.CurrentTheme
	b.InitControlBase(b, x, y, w, h)
	b.InitMouseInteraction(b)
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, theme.Text, true)

	b.OnAfterUpdate = func() {
		knobColor := theme.ScrollKnob
		if b.IsHovering || b.IsPressed {
			knobColor = theme.ScrollKnobHover
		}
		b.TextColor = knobColor
	}
}

// ScrollSliderVはスクロールバー用のスライド範囲
type ScrollSliderV struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.Grouping

	knob                InteractiveControl
	ViewRange, AllRange *int
	Value               float64
	OnSlide             func()
}

// ScrollSliderV生成
func NewSliderV(x, y, w, h int) *ScrollSliderV {
	s := &ScrollSliderV{}
	s.InitSliderV(x, y, w, h)
	return s
}

func (s *ScrollSliderV) InitSliderV(x, y, w, h int) {
	theme := parts.CurrentTheme
	s.InitControlBase(s, x, y, w, h)
	s.InitMouseInteraction(s)
	s.InitGrouping(s)

	s.AutoResizable = true

	s.knob.InitInteractiveControl(x, y, w, 60)

	// ノブのカスタム描画（角丸カプセル）
	s.knob.OnDraw = func(screen *ebiten.Image) {
		k := &s.knob
		gx, gy := k.GetGlobalPos()

		// 左右にパディングを入れる (inset)
		inset := 3
		width := k.Width - inset*2
		if width < 2 {
			width = 2
		}

		x := float32(gx + inset)
		y := float32(gy)
		w := float32(width)
		h := float32(k.Height)

		// 色決定
		knobColor := theme.ScrollKnob
		if k.IsHovering || k.IsDragging {
			knobColor = theme.ScrollKnobHover
		}

		// カプセル描画 (上下半円 + 中央矩形)
		r := w / 2

		vector.FillCircle(screen, x+r, y+r, r, knobColor, true)
		vector.FillCircle(screen, x+r, y+h-r, r, knobColor, true)
		vector.FillRect(screen, x, y+r, w, h-r*2, knobColor, false)
	}

	s.AddChild(&s.knob)

	dragOffsetY := 0
	s.knob.OnDragStart = func(x, y int) {
		_, gy := s.knob.GetGlobalPos()
		dragOffsetY = y - gy
	}
	s.knob.OnDrag = func(x, y int) {
		_, oy := s.GetGlobalPos()
		targetY := y - dragOffsetY - oy
		if s.Height-s.knob.Height > 0 {
			targetValue := float64(targetY) / float64(s.Height-s.knob.Height) * float64(*s.AllRange-*s.ViewRange)
			s.Move(targetValue - s.Value)
		}
	}

	s.OnRepeat = func() {
		// クリックした位置がツマミより上か下かで移動方向を決める
		_, y, _ := s.GetPosition()
		_, ky := s.knob.GetGlobalPos()
		if y < ky {
			s.Move(-float64(*s.ViewRange))
		} else if y > ky+s.knob.Height {
			s.Move(float64(*s.ViewRange))
		}
	}
}

// 指定した量だけスクロールする
func (s *ScrollSliderV) Move(delta float64) {
	s.Value = clampValue(s.Value+delta, *s.AllRange, *s.ViewRange)
	s.knob.Y = calculateKnobPos(s.Value, s.Height, s.knob.Height, *s.ViewRange, *s.AllRange)

	if s.OnSlide != nil {
		s.OnSlide()
	}
}

// ScrollSliderVの範囲設定
func (s *ScrollSliderV) Layout() {
	if *s.ViewRange >= *s.AllRange {
		s.knob.Height = s.Height
		s.knob.Y = 0
		s.Value = 0
		return
	}

	// まず Value をクランプ
	s.Value = clampValue(s.Value, *s.AllRange, *s.ViewRange)

	// 次にツマミサイズと位置を計算
	s.knob.Height = int(float64(s.Height) * float64(*s.ViewRange) / float64(*s.AllRange))
	if s.knob.Height < s.Width {
		s.knob.Height = s.Width
	}
	s.knob.Y = calculateKnobPos(s.Value, s.Height, s.knob.Height, *s.ViewRange, *s.AllRange)
}

// ScrollBarVは縦方向にスクロールするための複合Control
type ScrollBarV struct {
	parts.ControlBase
	parts.Drawable
	parts.Grouping

	buttonUp   ScrollButton
	slider     ScrollSliderV
	buttonDown ScrollButton

	OnSlide func()
}

// ScrollBarV生成
func NewScrollBarV(x, y, w, h int) *ScrollBarV {
	s := &ScrollBarV{}
	s.InitScrollBarV(x, y, w, h)
	return s
}

func (s *ScrollBarV) InitScrollBarV(x, y, w, h int) {
	s.InitControlBase(s, x, y, w, h)
	s.InitDrawable(s)
	s.InitGrouping(s)
	s.AutoLayout = parts.AutoLayoutFitV(0)

	s.buttonUp.InitScrollButton(0, 0, w, w, "▲", w/2)
	s.slider.InitSliderV(x, w, w, h)
	s.buttonDown.InitScrollButton(0, 0, w, w, "▼", w/2)

	s.AddChild(&s.buttonUp)
	s.AddChild(&s.slider)
	s.AddChild(&s.buttonDown)

	// ボタンの挙動設定
	// 押しっぱなしで連続スクロールするようにOnRepeatを使用
	s.buttonUp.OnRepeat = func() {
		s.slider.Move(-float64(*s.slider.ViewRange) * 0.1)
	}
	s.buttonDown.OnRepeat = func() {
		s.slider.Move(float64(*s.slider.ViewRange) * 0.1)
	}

	s.slider.OnSlide = func() {
		if s.OnSlide != nil {
			s.OnSlide()
		}
	}

	s.OnDraw = func(screen *ebiten.Image) {
		theme := parts.CurrentTheme

		gx, gy := s.GetGlobalPos()

		x := float32(gx)
		y := float32(gy)
		w := float32(s.Width)
		h := float32(s.Height)

		// カプセル描画 (上下半円 + 中央矩形)
		r := w / 2

		vector.FillCircle(screen, x+r, y+r, r, theme.ScrollBackground, true)
		vector.FillCircle(screen, x+r, y+h-r, r, theme.ScrollBackground, true)
		vector.FillRect(screen, x, y+r, w, h-r*2, theme.ScrollBackground, false)
	}
}

func (s *ScrollBarV) SetViewRange(viewrange *int) {
	s.slider.ViewRange = viewrange
}

func (s *ScrollBarV) SetMaxRange(allrange *int) {
	s.slider.AllRange = allrange
}

func (s *ScrollBarV) GetValue() float64 {
	return s.slider.Value
}

// ScrollSliderHは横方向のスクロールバー用のスライド範囲
type ScrollSliderH struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.Grouping

	knob                InteractiveControl
	ViewRange, AllRange *int
	Value               float64
	OnSlide             func()
}

// ScrollSliderH生成
func NewSliderH(x, y, w, h int) *ScrollSliderH {
	s := &ScrollSliderH{}
	s.InitSliderH(x, y, w, h)
	return s
}

func (s *ScrollSliderH) InitSliderH(x, y, w, h int) {
	theme := parts.CurrentTheme
	s.InitControlBase(s, x, y, w, h)
	s.InitMouseInteraction(s)
	s.InitGrouping(s)
	s.AutoResizable = true

	s.knob.InitInteractiveControl(x, 0, 60, h)

	// ノブのカスタム描画
	s.knob.OnDraw = func(screen *ebiten.Image) {
		k := &s.knob
		gx, gy := k.GetGlobalPos()

		// 上下にパディング (inset)
		inset := 3
		height := k.Height - inset*2
		if height < 2 {
			height = 2
		}

		x := float32(gx)
		y := float32(gy + inset)
		w := float32(k.Width)
		h := float32(height)

		knobColor := theme.ScrollKnob
		if k.IsHovering || k.IsDragging {
			knobColor = theme.ScrollKnobHover
		}

		r := h / 2
		vector.FillCircle(screen, x+r, y+r, r, knobColor, true)
		vector.FillCircle(screen, x+w-r, y+r, r, knobColor, true)
		vector.FillRect(screen, x+r, y, w-r*2, h, knobColor, false)
	}

	s.AddChild(&s.knob)

	var dragOffsetX int
	s.knob.OnDragStart = func(x, y int) {
		gx, _ := s.knob.GetGlobalPos()
		dragOffsetX = x - gx
	}
	s.knob.OnDrag = func(x, y int) {
		ox, _ := s.GetGlobalPos()
		targetX := x - dragOffsetX - ox
		if s.Width-s.knob.Width > 0 {
			targetValue := float64(targetX) / float64(s.Width-s.knob.Width) * float64(*s.AllRange-*s.ViewRange)
			s.Move(targetValue - s.Value)
		}
	}

	s.OnRepeat = func() {
		// クリックした位置がツマミより左か右かで移動方向を決める
		x, _, _ := s.GetPosition()
		kx, _ := s.knob.GetGlobalPos()
		if x < kx {
			s.Move(-float64(*s.ViewRange))
		} else if x > kx+s.knob.Width {
			s.Move(float64(*s.ViewRange))
		}
	}
}

// 指定した量だけスクロールする
func (s *ScrollSliderH) Move(delta float64) {
	s.Value = clampValue(s.Value+delta, *s.AllRange, *s.ViewRange)
	s.knob.X = calculateKnobPos(s.Value, s.Width, s.knob.Width, *s.ViewRange, *s.AllRange)

	if s.OnSlide != nil {
		s.OnSlide()
	}
}

// ScrollSliderHの範囲設定
func (s *ScrollSliderH) Layout() {
	if *s.ViewRange >= *s.AllRange {
		s.knob.Width = s.Width
		s.knob.X = 0
		s.Value = 0
		return
	}

	// まず Value をクランプ
	s.Value = clampValue(s.Value, *s.AllRange, *s.ViewRange)

	// 次にツマミサイズと位置を計算
	s.knob.Width = int(float64(s.Width) * float64(*s.ViewRange) / float64(*s.AllRange))
	if s.knob.Width < s.Height {
		s.knob.Width = s.Height
	}
	s.knob.X = calculateKnobPos(s.Value, s.Width, s.knob.Width, *s.ViewRange, *s.AllRange)
}

// ScrollBarHは横方向にスクロールするための複合Control
type ScrollBarH struct {
	parts.ControlBase
	parts.Drawable
	parts.Grouping

	buttonLeft  ScrollButton
	Slider      ScrollSliderH
	buttonRight ScrollButton

	OnSlide func()
}

// ScrollBarH生成
func NewScrollBarH(x, y, w, h int) *ScrollBarH {
	s := &ScrollBarH{}
	s.InitScrollBarH(x, y, w, h)
	return s
}

func (s *ScrollBarH) InitScrollBarH(x, y, w, h int) {
	s.InitControlBase(s, x, y, w, h)
	s.InitDrawable(s)
	s.InitGrouping(s)
	s.AutoLayout = parts.AutoLayoutFitH(0)

	s.buttonLeft.InitScrollButton(0, 0, h, h, "◀", h/2)
	s.Slider.InitSliderH(h, 0, w, h)
	s.buttonRight.InitScrollButton(0, 0, h, h, "▶", h/2)

	s.AddChild(&s.buttonLeft)
	s.AddChild(&s.Slider)
	s.AddChild(&s.buttonRight)

	// ボタンの挙動設定
	// 押しっぱなしで連続スクロールするようにOnRepeatを使用
	s.buttonLeft.OnRepeat = func() {
		s.Slider.Move(-float64(*s.Slider.ViewRange) * 0.1)
	}
	s.buttonRight.OnRepeat = func() {
		s.Slider.Move(float64(*s.Slider.ViewRange) * 0.1)
	}

	s.Slider.OnSlide = func() {
		if s.OnSlide != nil {
			s.OnSlide()
		}
	}

	s.OnDraw = func(screen *ebiten.Image) {
		theme := parts.CurrentTheme

		gx, gy := s.GetGlobalPos()

		x := float32(gx)
		y := float32(gy)
		w := float32(s.Width)
		h := float32(s.Height)

		// カプセル描画 (上下半円 + 中央矩形)
		r := h / 2

		vector.FillCircle(screen, x+r, y+r, r, theme.ScrollBackground, true)
		vector.FillCircle(screen, x+w-r, y+r, r, theme.ScrollBackground, true)
		vector.FillRect(screen, x+r, y, w-r*2, h, theme.ScrollBackground, false)
	}
}

func (s *ScrollBarH) SetViewRange(viewrange *int) {
	s.Slider.ViewRange = viewrange
}

func (s *ScrollBarH) SetMaxRange(allrange *int) {
	s.Slider.AllRange = allrange
}

func (s *ScrollBarH) GetValue() float64 {
	return s.Slider.Value
}
