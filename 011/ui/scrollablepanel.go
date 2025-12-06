package ui

import "MyProject/ui/parts"

// ScrollablePanelはスクロール可能なパネル
type ScrollablePanel struct {
	BlankWidget

	topGroup   BlankWidget
	Area       BlankWidget
	Panel      BlankWidget
	ScrollbarV ScrollBarV

	bottomGroup BlankWidget
	ScrollbarH  ScrollBarH
	corner      BlankWidget

	OnSlide func(x, y float64)
}

// ScrollablePanel生成
func NewScrollablePanel(x, y, w, h, sbw int) *ScrollablePanel {
	sp := &ScrollablePanel{}
	sp.InitScrollablePanel(x, y, w, h, sbw)
	return sp
}

func (sp *ScrollablePanel) InitScrollablePanel(x, y, w, h, sbw int) {
	sp.InitBlankWidget(x, y, w, h)
	sp.AutoLayout = parts.AutoLayoutFitV

	// 上部（パネル＋縦スクロールバー）
	sp.topGroup.InitBlankWidget(0, 0, w, h-sbw)
	sp.Grouping.AddChild(&sp.topGroup)
	sp.topGroup.AutoResizable = true
	sp.topGroup.AutoLayout = parts.AutoLayoutFitH
	sp.topGroup.ClippingFlag = false

	sp.Area.InitBlankWidget(0, 0, 0, 0)
	sp.topGroup.AddChild(&sp.Area)
	sp.Area.AutoResizable = true
	sp.ScrollbarV.InitScrollBarV(0, 0, sbw, 0)
	sp.topGroup.AddChild(&sp.ScrollbarV)

	sp.Panel.InitBlankWidget(0, 0, 0, 0)
	sp.Area.AddChild(&sp.Panel)
	sp.Panel.ClippingFlag = false

	// 下部（横スクロールバー＋コーナー）
	sp.bottomGroup.InitBlankWidget(0, 0, w, sbw)
	sp.Grouping.AddChild(&sp.bottomGroup)
	sp.bottomGroup.AutoLayout = parts.AutoLayoutFitH
	sp.bottomGroup.ClippingFlag = false

	sp.ScrollbarH.InitScrollBarH(0, 0, 0, sbw)
	sp.bottomGroup.AddChild(&sp.ScrollbarH)
	sp.ScrollbarH.AutoResizable = true
	sp.corner.InitBlankWidget(0, 0, sbw, sbw) // コーナーの空白
	sp.bottomGroup.AddChild(&sp.corner)

	// スクロールバーのスライド時の動作
	sp.ScrollbarV.OnSlide = func() {
		if sp.OnSlide != nil {
			sp.OnSlide(sp.ScrollbarH.GetValue(), sp.ScrollbarV.GetValue())
		}
	}
	sp.ScrollbarH.OnSlide = func() {
		if sp.OnSlide != nil {
			sp.OnSlide(sp.ScrollbarH.GetValue(), sp.ScrollbarV.GetValue())
		}
	}
	sp.OnLayout = sp.OnLayoutFunction
	sp.OnSlide = sp.OnSlideFunction
}

// ScrollablePanelに対してのAddChildはPanelに委譲する
func (sp *ScrollablePanel) AddChild(c parts.Widget) {
	sp.Panel.AddChild(c)
}

// ScrollablePanelのレイアウト
func (sp *ScrollablePanel) OnLayoutFunction() {
	// 縦スクロールバーの表示制御
	oldV := sp.ScrollbarV.Visible
	if sp.ScrollbarV.slider.ViewRange >= sp.ScrollbarV.slider.AllRange {
		sp.ScrollbarV.Visible = false
		sp.corner.Visible = false
	} else {
		sp.ScrollbarV.Visible = true
		sp.corner.Visible = true
	}

	// 横スクロールバーの表示制御
	oldH := sp.ScrollbarH.Visible
	if sp.ScrollbarH.slider.ViewRange >= sp.ScrollbarH.slider.AllRange {
		sp.ScrollbarH.Visible = false
		sp.bottomGroup.Visible = false
	} else {
		sp.ScrollbarH.Visible = true
		sp.bottomGroup.Visible = true
	}

	// レイアウトの更新
	if oldV != sp.ScrollbarV.Visible || oldH != sp.ScrollbarH.Visible {
		sp.Grouping.Layout()
	}

	// Panel配下の子コントロールの描画範囲を算出
	maxX, maxY := 0, 0
	for _, ch := range sp.Panel.Children {
		cb := ch.GetWidgetBase()
		maxX = max(cb.X+cb.Width, maxX)
		maxY = max(cb.Y+cb.Height, maxY)
	}
	sp.ScrollbarH.SetRange(float64(sp.Area.Width), float64(maxX))
	sp.ScrollbarV.SetRange(float64(sp.Area.Height), float64(maxY))
	sp.Panel.X, sp.Panel.Y = -int(sp.ScrollbarH.GetValue()), -int(sp.ScrollbarV.GetValue())

}

func (sp *ScrollablePanel) OnSlideFunction(x, y float64) {
	sp.Panel.X, sp.Panel.Y = -int(x), -int(y)
}
