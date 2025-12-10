package ui

import (
	"MyProject/ui/parts"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

var PopupContainer *popupContainer

type popupContainer struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.Grouping
	parts.Updatable
}

func NewPopupContainer() *popupContainer {
	p := &popupContainer{}
	p.InitPopupContainer()
	return p
}

func (p *popupContainer) InitPopupContainer() {
	p.InitControlBase(p, 0, 0, 0, 0)
	p.InitMouseInteraction(p)
	p.InitGrouping(p)
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

func (p *popupContainer) Close() {
	p.Children = []parts.Control{}
	p.Visible = false
}
