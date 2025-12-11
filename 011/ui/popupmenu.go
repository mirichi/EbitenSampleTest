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
}

// PopupMenuはドロップダウンやコンテキストメニューに使用するポップアップメニュー
type PopupMenu struct {
	parts.ControlBase
	parts.Drawable
	parts.Grouping

	Items           []MenuItem
	OnClose         func()
	ItemHeight      int
	SeparatorHeight int
	Margin          int
}

// PopupMenu生成
func NewPopupMenu(x, y int, items []MenuItem) *PopupMenu {
	pm := &PopupMenu{}
	pm.InitPopupMenu(x, y, items)
	return pm
}

func (pm *PopupMenu) InitPopupMenu(x, y int, items []MenuItem) {
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
	for i, menuItem := range items {
		if menuItem.IsSeparator {
			sep := newPopupMenuSeparator(width-pm.Margin*2, pm.SeparatorHeight)
			sep.Y = currentY
			pm.AddChild(sep)
			currentY += pm.SeparatorHeight
		} else {
			item := newPopupMenuItem(pm, i, menuItem, width-pm.Margin*2, pm.ItemHeight)
			item.Y = currentY
			pm.AddChild(item)
			currentY += pm.ItemHeight
		}
	}
}

// Close はメニューを閉じる
func (pm *PopupMenu) Close() {
	PopupContainer.Close()
}

// --- PopupMenuSeparator ---

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

// --- PopupMenuItem ---

// popupMenuItemはPopupMenuの各項目
type popupMenuItem struct {
	InteractiveControl
	parts.TextDrawable

	menu     *PopupMenu
	index    int
	action   func()
	disabled bool
}

func newPopupMenuItem(menu *PopupMenu, index int, menuItem MenuItem, width, height int) *popupMenuItem {
	theme := parts.CurrentTheme
	item := &popupMenuItem{}
	item.menu = menu
	item.index = index
	item.action = menuItem.Action
	item.disabled = menuItem.Disabled

	item.InitInteractiveControl(0, 0, width, height)

	// Disabledならグレーテキスト
	textColor := theme.Text
	if menuItem.Disabled {
		textColor = theme.DisabledText
	}
	item.InitTextDrawable(item, menuItem.Text, 14, parts.AlignLeft, parts.AlignCenter, 8, 0, textColor, true)

	// 描画
	item.OnDraw = func(screen *ebiten.Image) {
		gx, gy := item.GetGlobalPos()
		if item.IsMouseOver && !item.disabled {
			vector.FillRect(screen, float32(gx+menu.Margin), float32(gy), float32(item.Width), float32(item.Height), theme.PopupHover, false)
		}
	}

	// クリック
	item.OnRelease = func() {
		if item.disabled {
			return
		}
		if item.action != nil {
			item.action()
		}
		menu.Close()
	}

	return item
}
