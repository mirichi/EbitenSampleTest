package ui

import (
	"MyProject/ui/parts"
)

// ScrollablePanelはスクロール可能なパネル
type ScrollablePanel struct {
	GroupingControl

	topGroup   GroupingControl
	Area       GroupingControl
	Panel      GroupingControl
	ScrollbarV ScrollBarV

	bottomGroup GroupingControl
	ScrollbarH  ScrollBarH
	corner      BlankControl
}

// ScrollablePanel生成
func NewScrollablePanel(x, y, sbw int) *ScrollablePanel {
	sp := &ScrollablePanel{}
	sp.InitScrollablePanel(x, y, sbw)
	return sp
}

func (sp *ScrollablePanel) InitScrollablePanel(x, y, sbw int) {
	sp.InitGroupingControl(x, y, 0, 0)
	sp.AutoLayout = parts.FlexLayoutV(parts.FlexStart, parts.FlexStretch, 0)

	// 上部（パネル＋縦スクロールバー）
	sp.topGroup.InitGroupingControl(0, 0, 0, 0)
	sp.Grouping.AddChild(&sp.topGroup)
	sp.topGroup.FlexGrow = 1
	sp.topGroup.AutoLayout = parts.FlexLayoutH(parts.FlexStart, parts.FlexStretch, 0)

	sp.Area.InitGroupingControl(0, 0, 0, 0)
	sp.topGroup.AddChild(&sp.Area)
	sp.Area.FlexGrow = 1
	sp.ScrollbarV.InitScrollBarV(0, 0, sbw, 0)
	sp.topGroup.AddChild(&sp.ScrollbarV)
	sp.Area.ClippingFlag = true

	sp.Panel.InitGroupingControl(0, 0, 0, 0)
	sp.Area.AddChild(&sp.Panel)

	// 下部（横スクロールバー＋コーナー）
	sp.bottomGroup.InitGroupingControl(0, 0, 0, sbw)
	sp.Grouping.AddChild(&sp.bottomGroup)
	sp.bottomGroup.AutoLayout = parts.FlexLayoutH(parts.FlexStart, parts.FlexStretch, 0)

	sp.ScrollbarH.InitScrollBarH(0, 0, 0, sbw)
	sp.ScrollbarH.FlexGrow = 1
	sp.bottomGroup.AddChild(&sp.ScrollbarH)
	sp.corner.InitBlankControl(0, 0, sbw, sbw) // コーナーの空白
	sp.bottomGroup.AddChild(&sp.corner)

	// スクロールバーのスライド時の動作
	sp.ScrollbarV.OnSlide = func() {
		sp.Panel.X, sp.Panel.Y = -int(sp.ScrollbarH.GetValue()), -int(sp.ScrollbarV.GetValue())
	}
	sp.ScrollbarH.OnSlide = func() {
		sp.Panel.X, sp.Panel.Y = -int(sp.ScrollbarH.GetValue()), -int(sp.ScrollbarV.GetValue())
	}
	sp.OnLayout = sp.OnLayoutFunction

	sp.SetViewRange(&sp.Area.Width, &sp.Area.Height)
}

// ScrollablePanelに対してのAddChildはPanelに委譲する
func (sp *ScrollablePanel) AddChild(c parts.Entity) {
	sp.Panel.AddChild(c)
}

// ScrollablePanelのレイアウト
func (sp *ScrollablePanel) OnLayoutFunction() {
	// 縦スクロールバーの表示制御
	oldV := sp.ScrollbarV.Visible
	if *sp.ScrollbarV.slider.ViewRange >= *sp.ScrollbarV.slider.AllRange {
		sp.ScrollbarV.Visible = false
		sp.corner.Visible = false
	} else {
		sp.ScrollbarV.Visible = true
		sp.corner.Visible = true
	}

	// 横スクロールバーの表示制御
	oldH := sp.ScrollbarH.Visible
	if *sp.ScrollbarH.Slider.ViewRange >= *sp.ScrollbarH.Slider.AllRange {
		sp.bottomGroup.Visible = false
	} else {
		sp.bottomGroup.Visible = true
	}

	// レイアウト実行中に表示制御が変わったら再レイアウト
	// 再帰呼び出しにならないように気を付ける
	if oldV != sp.ScrollbarV.Visible || oldH != sp.ScrollbarH.Visible {
		sp.Layout()
	}

	sp.Panel.X, sp.Panel.Y = -int(sp.ScrollbarH.GetValue()), -int(sp.ScrollbarV.GetValue())
}

func (sp *ScrollablePanel) SetMaxRange(allrangeH, allrangeV *int) {
	sp.ScrollbarH.SetMaxRange(allrangeH)
	sp.ScrollbarV.SetMaxRange(allrangeV)
}

func (sp *ScrollablePanel) SetViewRange(viewrangeH, viewrangeV *int) {
	sp.ScrollbarH.SetViewRange(viewrangeH)
	sp.ScrollbarV.SetViewRange(viewrangeV)
}
