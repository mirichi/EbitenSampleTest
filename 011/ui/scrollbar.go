package ui

import (
	"MyProject/ui/parts"
	"image/color"
)

// ScrollButtonはスクロールバー用のフォーカスを持たないボタン
type ScrollButton struct {
	InteractiveControl
	parts.TextDrawable
}

// ScrollButton生成
func NewScrollButton(x, y, w, h int, text string, size int) *ScrollButton {
	b := &ScrollButton{}
	b.InitScrollButton(x, y, w, h, text, size)
	return b
}

func (b *ScrollButton) InitScrollButton(x, y, w, h int, text string, size int) {
	b.InitInteractiveControl(x, y, w, h)
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
}

// ScrollSliderVはスクロールバー用のスライド範囲
type ScrollSliderV struct {
	InteractiveControl
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
	s.InitInteractiveControl(x, y, w, h)
	s.InitGrouping(s)

	s.BackColor = color.RGBA{0x20, 0x20, 0x20, 0xff}
	s.ClippingFlag = false
	s.AutoResizable = true

	s.knob.InitInteractiveControl(x, y, w, 60)
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
	s.Value += delta
	if s.Value < 0 {
		s.Value = 0
	}
	if s.Value > float64(*s.AllRange-*s.ViewRange) {
		s.Value = float64(*s.AllRange - *s.ViewRange)
	}

	// ツマミの位置を更新
	if *s.AllRange > *s.ViewRange {
		s.knob.Y = int(float64(s.Value) * float64(s.Height-s.knob.Height) / float64(*s.AllRange-*s.ViewRange))
	} else {
		s.knob.Y = 0
	}

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

	s.knob.Height = int(float64(s.Height) * float64(*s.ViewRange) / float64(*s.AllRange))
	s.knob.Y = int(float64(s.Value) * float64(s.Height-s.knob.Height) / float64(*s.AllRange-*s.ViewRange))

	if s.knob.Y < 0 {
		s.knob.Y = 0
		s.Value = 0
	}
	if s.knob.Y > s.Height-s.knob.Height {
		s.knob.Y = s.Height - s.knob.Height
		s.Value = float64(*s.AllRange - *s.ViewRange)
	}
}

// ScrollBarVは縦方向にスクロールするための複合Widget
type ScrollBarV struct {
	GroupingControl

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
	s.InitGroupingControl(x, y, w, h)
	s.ClippingFlag = false
	s.AutoLayout = parts.AutoLayoutFitV

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
	InteractiveControl
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
	s.InitInteractiveControl(x, y, w, h)
	s.InitGrouping(s)

	s.BackColor = color.RGBA{0x20, 0x20, 0x20, 0xff}
	s.ClippingFlag = false
	s.AutoResizable = true

	s.knob.InitInteractiveControl(x, 0, 60, h)
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
	s.Value += delta
	if s.Value < 0 {
		s.Value = 0
	}
	if s.Value > float64(*s.AllRange-*s.ViewRange) {
		s.Value = float64(*s.AllRange - *s.ViewRange)
	}

	// ツマミの位置を更新
	if *s.AllRange > *s.ViewRange {
		s.knob.X = int(float64(s.Value) * float64(s.Width-s.knob.Width) / float64(*s.AllRange-*s.ViewRange))
	} else {
		s.knob.X = 0
	}

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

	s.knob.Width = int(float64(s.Width) * float64(*s.ViewRange) / float64(*s.AllRange))
	s.knob.X = int(float64(s.Value) * float64(s.Width-s.knob.Width) / float64(*s.AllRange-*s.ViewRange))

	if s.knob.X < 0 {
		s.knob.X = 0
		s.Value = 0
	}
	if s.knob.X > s.Width-s.knob.Width {
		s.knob.X = s.Width - s.knob.Width
		s.Value = float64(*s.AllRange - *s.ViewRange)
	}
}

// ScrollBarHは横方向にスクロールするための複合Widget
type ScrollBarH struct {
	GroupingControl

	buttonLeft  ScrollButton
	slider      ScrollSliderH
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
	s.InitGroupingControl(x, y, w, h)
	s.ClippingFlag = false
	s.AutoLayout = parts.AutoLayoutFitH

	s.buttonLeft.InitScrollButton(0, 0, h, h, "◀", h/2)
	s.slider.InitSliderH(h, 0, w, h)
	s.buttonRight.InitScrollButton(0, 0, h, h, "▶", h/2)

	s.AddChild(&s.buttonLeft)
	s.AddChild(&s.slider)
	s.AddChild(&s.buttonRight)

	// ボタンの挙動設定
	// 押しっぱなしで連続スクロールするようにOnRepeatを使用
	s.buttonLeft.OnRepeat = func() {
		s.slider.Move(-float64(*s.slider.ViewRange) * 0.1)
	}
	s.buttonRight.OnRepeat = func() {
		s.slider.Move(float64(*s.slider.ViewRange) * 0.1)
	}

	s.slider.OnSlide = func() {
		if s.OnSlide != nil {
			s.OnSlide()
		}
	}
}

func (s *ScrollBarH) SetViewRange(viewrange *int) {
	s.slider.ViewRange = viewrange
}

func (s *ScrollBarH) SetMaxRange(allrange *int) {
	s.slider.AllRange = allrange
}

func (s *ScrollBarH) GetValue() float64 {
	return s.slider.Value
}
