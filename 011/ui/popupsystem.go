package ui

import (
	"MyProject/ui/parts"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

var PopupWidgets *PopupContainer

type PopupContainer struct {
	parts.MouseInteraction
	GroupingWidget
	parts.Updatable
}

func NewPopupContainer() *PopupContainer {
	p := &PopupContainer{}
	p.InitPopupContainer()
	return p
}

func (p *PopupContainer) InitPopupContainer() {
	p.InitMouseInteraction(p)
	p.InitGroupingWidget(0, 0, 0, 0)
	p.InitUpdatable(p)
	p.ClippingFlag = false
	p.Visible = false

	p.OnUpdate = func() {
		if runtime.GOOS == "js" {
			p.Width = 640
			p.Height = 480
		} else {
			p.Width, p.Height = ebiten.WindowSize()
		}
		if len(p.Children) > 0 {
			p.Visible = true
		}
	}

	p.OnPress = func() {
		p.Close()
	}
}

func (p *PopupContainer) Close() {
	p.Children = []parts.Widget{}
	p.Visible = false
}
