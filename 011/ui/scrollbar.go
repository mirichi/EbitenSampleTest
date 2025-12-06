package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ScrollButtonはスクロールバー用のフォーカスを持たないボタン
type ScrollButton struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.Drawable
	parts.TextDrawable
}

// ScrollButton生成
func NewScrollButton(x, y, w, h int, text string, size int) *ScrollButton {
	b := &ScrollButton{}
	b.InitScrollButton(nil, x, y, w, h, text, size)
	return b
}

func (b *ScrollButton) InitScrollButton(g parts.AddChilder, x, y, w, h int, text string, size int) {
	b.InitControlBase(b, x, y, w, h)
	b.InitMouseInteraction(b)
	b.InitDrawable(b)
	b.OnDraw = b.drawScrollButton
	b.InitTextDrawable(b, text, size, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
	if g != nil {
		g.AddChild(b)
	}
}

func (b *ScrollButton) drawScrollButton(screen *ebiten.Image) {
	gx, gy := b.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(b.Width), float32(b.Height), color.RGBA{0x60, 0x60, 0x60, 0xff}, false)
}

// ScrollKnobはスクロールバー用のツマミ
type ScrollKnob struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.Drawable
}

// ScrollKnob生成
func NewKnob(x, y, w, h int) *ScrollKnob {
	k := &ScrollKnob{}
	k.InitScrollKnob(nil, x, y, w, h)
	return k
}

func (k *ScrollKnob) InitScrollKnob(g parts.AddChilder, x, y, w, h int) {
	k.InitControlBase(k, x, y, w, h)
	k.InitMouseInteraction(k)
	k.InitDrawable(k)
	k.OnDraw = k.drawKnob
	if g != nil {
		g.AddChild(k)
	}
}

func (k *ScrollKnob) drawKnob(screen *ebiten.Image) {
	gx, gy := k.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(k.Width), float32(k.Height), color.RGBA{0x60, 0x60, 0x60, 0xff}, false)
}

// ScrollSliderVはスクロールバー用のスライド範囲
type ScrollSliderV struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.Drawable
	parts.Grouping

	knob                       ScrollKnob
	ViewRange, AllRange, Value float64
	OnSlide                    func()
}

// ScrollSliderV生成
func NewSliderV(x, y, w, h int) *ScrollSliderV {
	s := &ScrollSliderV{}
	s.InitSliderV(nil, x, y, w, h)
	return s
}

func (s *ScrollSliderV) InitSliderV(g parts.AddChilder, x, y, w, h int) {
	s.InitControlBase(s, x, y, w, h)
	s.InitMouseInteraction(s)
	s.InitDrawable(s)
	s.OnDraw = s.drawSliderV
	s.InitGrouping(s)
	s.ClippingFlag = false
	s.knob.InitScrollKnob(s, 0, y, w, 60)
	s.AutoResizable = true
	if g != nil {
		g.AddChild(s)
	}

	var dragOffsetY int
	s.knob.OnDragStart = func(x, y int) {
		_, gy := s.knob.GetGlobalPos()
		dragOffsetY = y - gy
	}
	s.knob.OnDrag = func(x, y int) {
		_, oy := s.GetGlobalPos()
		targetY := y - dragOffsetY - oy
		if s.Height-s.knob.Height > 0 {
			targetValue := float64(targetY) / float64(s.Height-s.knob.Height) * (s.AllRange - s.ViewRange)
			s.Move(targetValue - s.Value)
		}
	}

	s.OnRepeat = func() {
		// クリックした位置がツマミより上か下かで移動方向を決める
		_, y, _ := s.GetPosition()
		_, ky := s.knob.GetGlobalPos()
		if y < ky {
			s.Move(-s.ViewRange)
		} else if y > ky+s.knob.Height {
			s.Move(s.ViewRange)
		}
	}
}

// 指定した量だけスクロールする
func (s *ScrollSliderV) Move(delta float64) {
	s.Value += delta
	if s.Value < 0 {
		s.Value = 0
	}
	if s.Value > s.AllRange-s.ViewRange {
		s.Value = s.AllRange - s.ViewRange
	}

	// ツマミの位置を更新
	// SetRangeのロジックを流用したいが、SetRangeは引数が必要なのでここで再計算
	if s.AllRange > s.ViewRange {
		s.knob.Y = int(float64(s.Value) * float64(s.Height-s.knob.Height) / float64(s.AllRange-s.ViewRange))
	} else {
		s.knob.Y = 0
	}

	if s.OnSlide != nil {
		s.OnSlide()
	}
}

func (s *ScrollSliderV) drawSliderV(screen *ebiten.Image) {
	gx, gy := s.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(s.Width), float32(s.Height), color.RGBA{0x20, 0x20, 0x20, 0xff}, false)
}

// ScrollSliderVの範囲設定
func (s *ScrollSliderV) SetRange(viewrange, allrange float64) {
	s.ViewRange = viewrange
	s.AllRange = allrange
	if viewrange >= allrange {
		s.knob.Height = s.Height
		s.knob.Y = 0
		s.Value = 0
		return
	}
	s.knob.Height = int(float64(s.Height) * float64(s.ViewRange) / float64(s.AllRange))
	s.knob.Y = int(float64(s.Value) * float64(s.Height-s.knob.Height) / float64(s.AllRange-s.ViewRange))

	if s.knob.Y < 0 {
		s.knob.Y = 0
		s.Value = float64(s.knob.Y) / float64(s.Height-s.knob.Height) * (s.AllRange - s.ViewRange)
	}
	if s.knob.Y > s.Height-s.knob.Height {
		s.knob.Y = s.Height - s.knob.Height
		s.Value = float64(s.knob.Y) / float64(s.Height-s.knob.Height) * (s.AllRange - s.ViewRange)
	}
}

// ScrollBarVは縦方向にスクロールするための複合コントロール
type ScrollBarV struct {
	parts.ControlBase
	parts.Grouping

	buttonUp   ScrollButton
	slider     ScrollSliderV
	buttonDown ScrollButton

	OnSlide func()
}

// ScrollBarV生成
func NewScrollBarV(x, y, w, h int) *ScrollBarV {
	s := &ScrollBarV{}
	s.InitScrollBarV(nil, x, y, w, h)
	return s
}

func (s *ScrollBarV) InitScrollBarV(g parts.AddChilder, x, y, w, h int) {
	s.InitControlBase(s, x, y, w, h)
	s.InitGrouping(s)
	s.ClippingFlag = false
	s.buttonUp.InitScrollButton(s, 0, 0, w, w, "▲", w/2)
	s.slider.InitSliderV(s, x, w, w, h)
	s.buttonDown.InitScrollButton(s, 0, 0, w, w, "▼", w/2)
	s.Grouping.AutoLayout = parts.AutoLayoutFitV
	if g != nil {
		g.AddChild(s)
	}

	// ボタンの挙動設定
	// 押しっぱなしで連続スクロールするようにOnRepeatを使用
	s.buttonUp.OnRepeat = func() {
		s.slider.Move(-s.slider.ViewRange * 0.1)
	}
	s.buttonDown.OnRepeat = func() {
		s.slider.Move(s.slider.ViewRange * 0.1)
	}

	s.slider.OnSlide = func() {
		if s.OnSlide != nil {
			s.OnSlide()
		}
	}
}

func (s *ScrollBarV) SetRange(viewrange, allrange float64) {
	s.slider.SetRange(viewrange, allrange)
}

func (s *ScrollBarV) GetValue() float64 {
	return s.slider.Value
}

// ScrollSliderHは横方向のスクロールバー用のスライド範囲
type ScrollSliderH struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.Drawable
	parts.Grouping

	knob                       ScrollKnob
	ViewRange, AllRange, Value float64
	OnSlide                    func()
}

// ScrollSliderH生成
func NewSliderH(x, y, w, h int) *ScrollSliderH {
	s := &ScrollSliderH{}
	s.InitSliderH(nil, x, y, w, h)
	return s
}

func (s *ScrollSliderH) InitSliderH(g parts.AddChilder, x, y, w, h int) {
	s.InitControlBase(s, x, y, w, h)
	s.InitMouseInteraction(s)
	s.InitDrawable(s)
	s.OnDraw = s.drawSliderH
	s.InitGrouping(s)
	s.ClippingFlag = false
	s.knob.InitScrollKnob(s, x, 0, 60, h)
	s.AutoResizable = true
	if g != nil {
		g.AddChild(s)
	}

	var dragOffsetX int
	s.knob.OnDragStart = func(x, y int) {
		gx, _ := s.knob.GetGlobalPos()
		dragOffsetX = x - gx
	}
	s.knob.OnDrag = func(x, y int) {
		ox, _ := s.GetGlobalPos()
		targetX := x - dragOffsetX - ox
		if s.Width-s.knob.Width > 0 {
			targetValue := float64(targetX) / float64(s.Width-s.knob.Width) * (s.AllRange - s.ViewRange)
			s.Move(targetValue - s.Value)
		}
	}

	s.OnRepeat = func() {
		// クリックした位置がツマミより左か右かで移動方向を決める
		x, _, _ := s.GetPosition()
		kx, _ := s.knob.GetGlobalPos()
		if x < kx {
			s.Move(-s.ViewRange)
		} else if x > kx+s.knob.Width {
			s.Move(s.ViewRange)
		}
	}
}

// 指定した量だけスクロールする
func (s *ScrollSliderH) Move(delta float64) {
	s.Value += delta
	if s.Value < 0 {
		s.Value = 0
	}
	if s.Value > s.AllRange-s.ViewRange {
		s.Value = s.AllRange - s.ViewRange
	}

	// ツマミの位置を更新
	// SetRangeのロジックを流用したいが、SetRangeは引数が必要なのでここで再計算
	if s.AllRange > s.ViewRange {
		s.knob.X = int(float64(s.Value) * float64(s.Width-s.knob.Width) / float64(s.AllRange-s.ViewRange))
	} else {
		s.knob.X = 0
	}

	if s.OnSlide != nil {
		s.OnSlide()
	}
}

func (s *ScrollSliderH) drawSliderH(screen *ebiten.Image) {
	gx, gy := s.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(s.Width), float32(s.Height), color.RGBA{0x20, 0x20, 0x20, 0xff}, false)
}

// ScrollSliderHの範囲設定
func (s *ScrollSliderH) SetRange(viewrange, allrange float64) {
	s.ViewRange = viewrange
	s.AllRange = allrange
	if viewrange >= allrange {
		s.knob.Width = s.Width
		s.knob.X = 0
		s.Value = 0
		return
	}
	s.knob.Width = int(float64(s.Width) * float64(s.ViewRange) / float64(s.AllRange))
	s.knob.X = int(float64(s.Value) * float64(s.Width-s.knob.Width) / float64(s.AllRange-s.ViewRange))

	if s.knob.X < 0 {
		s.knob.X = 0
		s.Value = float64(s.knob.X) / float64(s.Width-s.knob.Width) * (s.AllRange - s.ViewRange)
	}
	if s.knob.X > s.Width-s.knob.Width {
		s.knob.X = s.Width - s.knob.Width
		s.Value = float64(s.knob.X) / float64(s.Width-s.knob.Width) * (s.AllRange - s.ViewRange)
	}
}

// ScrollBarHは横方向にスクロールするための複合コントロール
type ScrollBarH struct {
	parts.ControlBase
	parts.Grouping

	buttonLeft  ScrollButton
	slider      ScrollSliderH
	buttonRight ScrollButton

	OnSlide func()
}

// ScrollBarH生成
func NewScrollBarH(x, y, w, h int) *ScrollBarH {
	s := &ScrollBarH{}
	s.InitScrollBarH(nil, x, y, w, h)
	return s
}

func (s *ScrollBarH) InitScrollBarH(g parts.AddChilder, x, y, w, h int) {
	s.InitControlBase(s, x, y, w, h)
	s.InitGrouping(s)
	s.ClippingFlag = false
	s.buttonLeft.InitScrollButton(s, 0, 0, h, h, "◀", h/2)
	s.slider.InitSliderH(s, h, 0, w, h)
	s.buttonRight.InitScrollButton(s, 0, 0, h, h, "▶", h/2)
	s.Grouping.AutoLayout = parts.AutoLayoutFitH
	if g != nil {
		g.AddChild(s)
	}

	// ボタンの挙動設定
	// 押しっぱなしで連続スクロールするようにOnRepeatを使用
	s.buttonLeft.OnRepeat = func() {
		s.slider.Move(-s.slider.ViewRange * 0.1)
	}
	s.buttonRight.OnRepeat = func() {
		s.slider.Move(s.slider.ViewRange * 0.1)
	}

	s.slider.OnSlide = func() {
		if s.OnSlide != nil {
			s.OnSlide()
		}
	}
}

func (s *ScrollBarH) SetRange(viewrange, allrange float64) {
	s.slider.SetRange(viewrange, allrange)
}

func (s *ScrollBarH) GetValue() float64 {
	return s.slider.Value
}
