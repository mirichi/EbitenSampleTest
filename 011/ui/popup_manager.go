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
		p.Visible = false
	}

	p.OnRightPress = func() {
		p.Visible = false
	}
}

// ShowMenu は指定されたメニューを表示する
func (p *PopupManager) ShowMenu(menu *PopupMenu) {
	// 既存のメニューを閉じる（単純化のため現在は1つのみ表示）
	p.Children = []parts.Control{}

	// コンテナに追加
	p.AddChild(menu)
	p.Visible = true

	menu.OnClose = func() {
		p.Visible = false
	}

	// サブメニュー表示のコールバック
	menu.OnShowSubMenu = func(submenu *PopupMenu) {
		// サブメニューの位置を親項目の右側に配置
		p.AddChild(submenu)
		submenu.OnClose = menu.OnClose
		submenu.OnShowSubMenu = menu.OnShowSubMenu // 再帰的に設定
	}
}
