package ui

import (
	"MyProject/ui/parts"
	"image/color"

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

// ScrollBarVはスクロールバー用のスライド範囲
type ScrollBarV struct {
	InteractiveControl
	parts.Grouping

	knob                InteractiveControl
	ViewRange, AllRange *int
	Value               float64
	OnSlide             func()
}

// ScrollBarV生成
func NewScrollBarV(x, y, w, h int) *ScrollBarV {
	s := &ScrollBarV{}
	s.InitScrollBarV(x, y, w, h)
	return s
}

func (s *ScrollBarV) InitScrollBarV(x, y, w, h int) {
	theme := parts.CurrentTheme
	s.InitInteractiveControl(x, y, w, h)
	s.InitGrouping(s)

	s.BackColor = color.Transparent // トラックは透明に
	s.ClippingFlag = false

	s.knob.InitInteractiveControl(x, y, w, 60)
	s.knob.BackColor = color.Transparent // ノブ自体のControl背景は透明にし、カスタム描画する

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
			knobColor = theme.ButtonHover // ホバー時は少し明るく（既存テーマ流用）
		}

		// カプセル描画 (上下半円 + 中央矩形)
		r := w / 2

		vector.FillCircle(screen, x+r, y+r, r, knobColor, true)
		vector.FillCircle(screen, x+r, y+h-r, r, knobColor, true)
		vector.FillRect(screen, x, y+r, w, h-r*2, knobColor, true)
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
func (s *ScrollBarV) Move(delta float64) {
	s.Value = clampValue(s.Value+delta, *s.AllRange, *s.ViewRange)
	s.knob.Y = calculateKnobPos(s.Value, s.Height, s.knob.Height, *s.ViewRange, *s.AllRange)

	if s.OnSlide != nil {
		s.OnSlide()
	}
}

// ScrollSliderVの範囲設定
func (s *ScrollBarV) Layout() {
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

func (s *ScrollBarV) SetViewRange(viewrange *int) {
	s.ViewRange = viewrange
}

func (s *ScrollBarV) SetMaxRange(allrange *int) {
	s.AllRange = allrange
}

// ScrollBarHは横方向のスクロールバー用のスライド範囲
type ScrollBarH struct {
	InteractiveControl
	parts.Grouping

	knob                InteractiveControl
	ViewRange, AllRange *int
	Value               float64
	OnSlide             func()
}

// ScrollBarH生成
func NewScrollBarH(x, y, w, h int) *ScrollBarH {
	s := &ScrollBarH{}
	s.InitScrollBarH(x, y, w, h)
	return s
}

func (s *ScrollBarH) InitScrollBarH(x, y, w, h int) {
	theme := parts.CurrentTheme
	s.InitInteractiveControl(x, y, w, h)
	s.InitGrouping(s)

	s.BackColor = color.Transparent
	s.ClippingFlag = false

	s.knob.InitInteractiveControl(x, 0, 60, h)
	s.knob.BackColor = color.Transparent

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
			knobColor = theme.ButtonHover
		}

		r := h / 2
		vector.FillCircle(screen, x+r, y+r, r, knobColor, true)
		vector.FillCircle(screen, x+w-r, y+r, r, knobColor, true)
		vector.FillRect(screen, x+r, y, w-r*2, h, knobColor, true)
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
func (s *ScrollBarH) Move(delta float64) {
	s.Value = clampValue(s.Value+delta, *s.AllRange, *s.ViewRange)
	s.knob.X = calculateKnobPos(s.Value, s.Width, s.knob.Width, *s.ViewRange, *s.AllRange)

	if s.OnSlide != nil {
		s.OnSlide()
	}
}

// ScrollSliderHの範囲設定
func (s *ScrollBarH) Layout() {
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

func (s *ScrollBarH) SetViewRange(viewrange *int) {
	s.ViewRange = viewrange
}

func (s *ScrollBarH) SetMaxRange(allrange *int) {
	s.AllRange = allrange
}
