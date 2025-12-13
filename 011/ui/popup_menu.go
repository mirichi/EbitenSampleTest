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

	Items           []*MenuItem
	OnCloseAll      func()
	ItemHeight      int
	SeparatorHeight int
	Margin          int
	SubMenu         *PopupMenu
}

// PopupMenu生成
func NewPopupMenu(x, y int, items []*MenuItem) *PopupMenu {
	pm := &PopupMenu{}
	pm.InitPopupMenu(x, y, items)
	return pm
}

func (pm *PopupMenu) InitPopupMenu(x, y int, items []*MenuItem) {
	pm.Items = items
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

	// 背景描画
	pm.OnDraw = func(screen *ebiten.Image) {
		theme := parts.CurrentTheme
		gx, gy := pm.GetGlobalPos()
		vector.FillRect(screen, float32(gx), float32(gy), float32(pm.Width), float32(pm.Height), theme.PopupBackground, false)
		vector.StrokeRect(screen, float32(gx), float32(gy), float32(pm.Width), float32(pm.Height), 1, theme.FocusBorder, false)
	}

	// メニュー項目を生成
	currentY := pm.Margin
	for _, menuItem := range items {
		if menuItem.IsSeparator {
			sep := newPopupMenuSeparator(width-pm.Margin*2, pm.SeparatorHeight)
			sep.Y = currentY
			pm.AddChild(sep)
			currentY += pm.SeparatorHeight
		} else {
			item := newPopupMenuItem(pm, menuItem, width-pm.Margin*2, pm.ItemHeight)
			item.Y = currentY
			pm.AddChild(item)
			currentY += pm.ItemHeight
		}
	}
}

// CloseSubMenuはサブメニューを閉じる
func (pm *PopupMenu) CloseSubMenu() {
	if pm.SubMenu != nil {
		pm.RemoveChild(pm.SubMenu)
		pm.SubMenu = nil
	}
}

// CloseAll はメニューをすべて閉じる
func (pm *PopupMenu) CloseAll() {
	if pm.OnCloseAll != nil {
		pm.OnCloseAll()
	}
}

// ShowSubMenuはサブメニューを表示する
func (pm *PopupMenu) ShowSubMenu(x, y int, submenu []*MenuItem) {
	pm.CloseSubMenu()
	subMenu := NewPopupMenu(x, y, submenu)
	subMenu.OnCloseAll = pm.OnCloseAll
	pm.AddChild(subMenu)
	pm.SubMenu = subMenu
}

// popupMenuSeparatorはセパレータ
type popupMenuSeparator struct {
	InteractiveControl
}

func newPopupMenuSeparator(width, height int) *popupMenuSeparator {
	sep := &popupMenuSeparator{}
	sep.InitInteractiveControl(0, 0, width, height)

	sep.OnDraw = func(screen *ebiten.Image) {
		theme := parts.CurrentTheme
		gx, gy := sep.GetGlobalPos()
		lineY := float32(gy) + float32(height)/2
		vector.StrokeLine(screen, float32(gx)+4, lineY, float32(gx+width)-4, lineY, 1, theme.DisabledText, false)
	}

	return sep
}

// popupMenuItemはPopupMenuの各項目
type popupMenuItem struct {
	InteractiveControl
	parts.TextDrawable

	Timer    parts.Timer
	menu     *PopupMenu
	menuItem *MenuItem
}

func newPopupMenuItem(menu *PopupMenu, menuItem *MenuItem, width, height int) *popupMenuItem {
	theme := parts.CurrentTheme
	item := &popupMenuItem{}
	item.menu = menu
	item.menuItem = menuItem

	item.InitInteractiveControl(0, 0, width, height)

	// Disabledならグレーテキスト
	textColor := theme.Text
	if menuItem.Disabled {
		textColor = theme.DisabledText
	}
	item.InitTextDrawable(item, menuItem.Text, 14, parts.AlignLeft, parts.AlignCenter, 8, 0, textColor, true)
	item.Timer.InitTimer(item, 30, false)

	// 描画
	item.OnDraw = func(screen *ebiten.Image) {
		// 選択色
		if item.IsMouseOver && !item.menuItem.Disabled {
			gx, gy := item.GetGlobalPos()
			vector.FillRect(screen, float32(gx+menu.Margin), float32(gy), float32(item.Width), float32(item.Height), theme.PopupHover, false)
		}
	}

	// クリック
	item.OnRelease = func() {
		if item.menuItem.Disabled {
			return
		}
		if item.menuItem.Action != nil {
			item.menuItem.Action()
			item.menu.CloseAll()
		}
		if item.menuItem.SubMenu != nil {
			item.menu.ShowSubMenu(item.X+item.Width, item.Y, item.menuItem.SubMenu)
			item.Timer.Stop()
		}
	}

	item.OnHoverStart = func() {
		if item.menuItem.SubMenu != nil {
			item.Timer.Start()
		}
	}

	item.OnHoverEnd = func() {
		item.Timer.Stop()
	}

	item.Timer.OnTimeout = func() {
		item.menu.ShowSubMenu(item.X+item.Width, item.Y, item.menuItem.SubMenu)
	}

	return item
}
