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

	p.OnBeforeHandleInput = func() {
		if runtime.GOOS == "js" {
			p.Width = 640
			p.Height = 480
		} else {
			p.Width, p.Height = ebiten.WindowSize()
		}
	}

	p.OnPress = func() {
		p.Visible = false
	}

	p.OnRightPress = func() {
		p.Visible = false
	}
}

// ShowMenu は指定されたメニューを表示する
func (p *PopupManager) ShowMenu(x, y int, menu []*MenuItem) {
	// 既存のメニューを閉じる（単純化のため現在は1つのみ表示）
	p.Children = []parts.Control{}

	popupMenu := NewPopupMenu(x, y, menu)

	// コンテナに追加
	p.AddChild(popupMenu)
	p.Visible = true

	popupMenu.OnCloseAll = func() {
		p.Visible = false

		// ループ中に消すとバグるのでUpdate後に1回だけ実行するようにして消す
		p.OnAfterUpdate = func() {
			p.Children = []parts.Control{}
			p.OnAfterUpdate = nil
		}
	}
}
