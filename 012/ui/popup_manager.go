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
	p.Visible = false

	p.AddBeforeUpdateFunction(func() {
		if runtime.GOOS == "js" {
			p.Width = 640
			p.Height = 480
		} else {
			w, h := ebiten.WindowSize()
			p.Width, p.Height = float64(w), float64(h)
		}
	})

	p.OnPress = func() {
		p.Visible = false
	}

	p.OnRightPress = func() {
		p.Visible = false
	}
}

// ShowMenu は指定されたメニューを表示する
func (p *PopupManager) ShowMenu(x, y float64, menu []*MenuItem) {
	// 既存のメニューを閉じる（単純化のため現在は1つのみ表示）
	p.Children = []parts.Entity{}

	popupMenu := NewPopupMenu(x, y, menu)

	// コンテナに追加
	p.AddChild(popupMenu)
	p.Visible = true

	closeflg := false
	p.AddAfterUpdateFunction(func() {
		if closeflg {
			p.Children = []parts.Entity{}
			closeflg = false
		}
	})
	popupMenu.OnCloseAll = func() {
		p.Visible = false
		closeflg = true
	}
}
