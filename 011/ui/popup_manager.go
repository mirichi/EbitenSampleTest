package ui

import (
	"MyProject/ui/parts"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

type PopupManager struct {
	parts.ControlBase
	parts.MouseInteraction
	parts.Grouping
}

func NewPopupManager() *PopupManager {
	p := &PopupManager{}
	p.InitPopupManager()
	return p
}

func (p *PopupManager) InitPopupManager() {
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

func (p *PopupManager) Close() {
	p.Children = []parts.Control{}
	p.Visible = false
}

// AddChildは子を追加し、Visibleをtrueにする
func (p *PopupManager) AddChild(c parts.Control) {
	p.Grouping.AddChild(c)
	p.Visible = true
}

// ShowMenu は指定されたメニューを表示する
func (p *PopupManager) ShowMenu(menu *PopupMenu) {
	// 既存のメニューを閉じる（単純化のため現在は1つのみ表示）
	p.Close()

	// コンテナに追加
	p.AddChild(menu)

	// Menu自身にもManagerへの参照を持たせると閉じる処理などがスムーズになるが
	// 現状はCallbackで対応
	menu.OnClose = func() {
		p.Close()
	}
}
