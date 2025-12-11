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
	p.ClippingFlag = false
	p.Visible = false

	p.OnBeforeHandleInput = func() {
		if runtime.GOOS == "js" {
			p.Width = 640
			p.Height = 480
		} else {
			p.Width, p.Height = ebiten.WindowSize()
		}
	}

	p.OnPress = func() {
		p.Close()
	}

	p.OnRightPress = func() {
		p.Close()
	}
}

func (p *popupContainer) Close() {
	p.Children = []parts.Control{}
	p.Visible = false
}

// AddChildは子を追加し、Visibleをtrueにする
func (p *popupContainer) AddChild(c parts.Control) {
	p.Grouping.AddChild(c)
	p.Visible = true
}
