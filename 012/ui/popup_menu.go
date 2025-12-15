package ui

import (
	"MyProject/ui/parts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// MenuItem はポップアップメニューの項目定義
type MenuItem struct {
	Text        string // 表示テキスト
	Action      func() // 選択時に実行される関数
	Disabled    bool   // trueなら選択不可・グレー表示
	IsSeparator bool   // trueなら区切り線として描画
	SubMenu     []*MenuItem
}

// PopupMenuはドロップダウンやコンテキストメニューに使用するポップアップメニュー
type PopupMenu struct {
	parts.ControlBase
	parts.Drawable
	parts.Grouping

	PopupItems      []*popupMenuItem
	OnCloseAll      func()
	ItemHeight      int
	SeparatorHeight int
	Margin          int
	SubMenu         *PopupMenu
	SubMenuParent   *popupMenuItem
	IsExpanded      bool
	Timer           parts.Timer
	hasHighlight    bool
}

// PopupMenu生成
func NewPopupMenu(x, y int, items []*MenuItem) *PopupMenu {
	pm := &PopupMenu{}
	pm.InitPopupMenu(x, y, items)
	return pm
}

func (pm *PopupMenu) InitPopupMenu(x, y int, items []*MenuItem) {
	pm.ItemHeight = 24
	pm.SeparatorHeight = 9
	pm.Margin = 2

	// サイズ計算（幅は固定、高さは項目に応じて）
	width := 150 + pm.Margin*2
	height := pm.Margin * 2
	for _, item := range items {
		if item.IsSeparator {
			height += pm.SeparatorHeight
		} else {
			height += pm.ItemHeight
		}
	}

	pm.InitControlBase(pm, x, y, width, height)
	pm.InitDrawable(pm)
	pm.InitGrouping(pm)
	pm.ClippingFlag = false
	// pm.AutoLayout = parts.AutoLayoutFitV(pm.Margin)
	pm.Timer.InitTimer(pm, 30, false)

	// 背景描画
	pm.OnDraw = func(screen *ebiten.Image) {
		theme := parts.CurrentTheme
		gx, gy := pm.GetGlobalPos()
		vector.FillRect(screen, float32(gx), float32(gy), float32(pm.Width), float32(pm.Height), theme.PopupBackground, false)
		vector.StrokeRect(screen, float32(gx), float32(gy), float32(pm.Width-1), float32(pm.Height-1), 1, theme.FocusBorder, false)
	}

	// メニュー項目を生成
	iy := pm.Margin
	for _, menuItem := range items {
		var item *popupMenuItem
		if menuItem.IsSeparator {
			item = newPopupMenuItem(pm.Margin, iy, pm.Width-pm.Margin*2, pm.SeparatorHeight, menuItem)
			iy = iy + pm.SeparatorHeight
		} else {
			item = newPopupMenuItem(pm.Margin, iy, pm.Width-pm.Margin*2, pm.ItemHeight, menuItem)
			iy = iy + pm.ItemHeight
		}

		item.OnPress = func() {
			if item.menuItem.SubMenu != nil {
				// ボタンを押した項目がサブメニューを持っていれば開く
				pm.ShowSubMenu(item.X+item.Width, item.Y, item.menuItem.SubMenu, item)
				pm.Timer.Stop()
			}
		}
		item.OnRelease = func() {
			if item.menuItem.Action != nil && !item.menuItem.Disabled {
				// ボタンを離した項目がアクションを持っていれば実行してメニューを閉じる
				item.menuItem.Action()
				pm.CloseAll()
			}

			if item.menuItem.SubMenu != nil && !item.menuItem.Disabled {
				// ボタンを離した項目がサブメニューを持っていれば開く
				pm.ShowSubMenu(item.X+item.Width, item.Y, item.menuItem.SubMenu, item)
				pm.Timer.Stop()
			}
		}

		pm.AddChild(item)
		pm.PopupItems = append(pm.PopupItems, item)
	}

	pm.OnAfterUpdate = func() {
		pm.hasHighlight = false
		activeItem := pm.findActiveItem()

		if activeItem != nil {
			// ハイライトされていない項目にカーソルを乗せたらタイマー開始
			if !activeItem.Highlight {
				pm.Timer.Resume()
			}
		} else {
			// 項目に乗っていなくてサブメニューが開いているときもタイマー開始
			if pm.SubMenu != nil {
				pm.Timer.Resume()
			}
		}

		// ハイライト解除
		for _, item := range pm.PopupItems {
			item.Highlight = false
		}

		// ハイライト設定
		if activeItem != nil {
			pm.IsExpanded = false
			pm.hasHighlight = true
			activeItem.Highlight = true
		}

		// サブメニューに乗っていたら親項目をハイライトする
		if pm.SubMenu != nil {
			if pm.SubMenu.hasHighlight && !pm.SubMenu.IsExpanded {
				pm.SubMenuParent.Highlight = true
				pm.hasHighlight = true
				pm.IsExpanded = true
				pm.Timer.Stop()
			} else if pm.IsExpanded {
				pm.hasHighlight = true
				pm.SubMenuParent.Highlight = true
				pm.Timer.Stop()
			}
		}
	}

	// Timerが発火したらサブメニューを開くか閉じる
	pm.Timer.OnTimer = func() {
		activeItem := pm.findActiveItem()
		if activeItem != nil {
			if activeItem.menuItem.SubMenu != nil {
				// カーソルに乗っている項目がサブメニューを持っている場合は開く
				pm.ShowSubMenu(activeItem.X+activeItem.Width, activeItem.Y, activeItem.menuItem.SubMenu, activeItem)
			} else {
				// サブメニューを持っていない場合は閉じる
				pm.CloseSubMenu()
			}
		} else {
			// 項目に乗っていない場合も閉じる
			pm.CloseSubMenu()
		}
	}
}

// カーソルが乗っている項目を返す
func (pm *PopupMenu) findActiveItem() *popupMenuItem {
	for _, item := range pm.PopupItems {
		if item.IsMouseOver && !item.menuItem.Disabled {
			return item
		}
	}
	return nil
}

// CloseSubMenuはサブメニューを閉じる
func (pm *PopupMenu) CloseSubMenu() {
	if pm.SubMenu != nil {
		pm.RemoveChild(pm.SubMenu)
		pm.SubMenu = nil
		pm.SubMenuParent = nil
	}
}

// CloseAllはメニューをすべて閉じる
func (pm *PopupMenu) CloseAll() {
	if pm.OnCloseAll != nil {
		pm.OnCloseAll()
	}
}

// ShowSubMenuはサブメニューを表示する
func (pm *PopupMenu) ShowSubMenu(x, y int, submenu []*MenuItem, item *popupMenuItem) {
	pm.CloseSubMenu()
	subMenu := NewPopupMenu(x, y, submenu)
	subMenu.OnCloseAll = pm.OnCloseAll
	pm.AddChild(subMenu)
	pm.SubMenu = subMenu
	pm.SubMenuParent = item
}

// popupMenuItemはPopupMenuの各項目
type popupMenuItem struct {
	InteractiveControl
	parts.TextDrawable

	menuItem      *MenuItem
	submenuMarker parts.TextDrawable
	Highlight     bool
}

func newPopupMenuItem(x, y, width, height int, menuItem *MenuItem) *popupMenuItem {
	theme := parts.CurrentTheme
	item := &popupMenuItem{}
	item.menuItem = menuItem

	item.InitInteractiveControl(x, y, width, height)

	if !item.menuItem.IsSeparator {
		// Disabledならグレーテキスト
		textColor := theme.Text
		if menuItem.Disabled {
			textColor = theme.DisabledText
		}

		item.InitTextDrawable(item, menuItem.Text, 14, parts.AlignLeft, parts.AlignCenter, 8, 0, textColor, true)
		if menuItem.SubMenu != nil {
			item.submenuMarker.InitTextDrawable(item, ">", 14, parts.AlignRight, parts.AlignCenter, -8, 0, textColor, true)
		}
	}

	// 描画
	item.OnDraw = func(screen *ebiten.Image) {
		if item.menuItem.IsSeparator {
			theme := parts.CurrentTheme
			gx, gy := item.GetGlobalPos()
			lineY := float32(gy) + float32(height)/2
			vector.StrokeLine(screen, float32(gx)+4, lineY, float32(gx+item.Width)-4, lineY, 1, theme.DisabledText, false)
		} else {
			if item.Highlight {
				gx, gy := item.GetGlobalPos()
				vector.FillRect(screen, float32(gx), float32(gy), float32(item.Width), float32(item.Height), theme.PopupHover, false)
			}
		}
	}

	return item
}
