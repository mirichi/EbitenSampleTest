package ui

import (
	"MyProject/ui/parts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PopupMenuはドロップダウンやコンテキストメニューに使用するポップアップメニュー
type PopupMenu struct {
	parts.ControlBase
	parts.Drawable
	parts.Grouping

	Items      []string
	OnSelect   func(index int)
	OnClose    func()
	ItemHeight int
}

// PopupMenu生成
func NewPopupMenu(x, y int, items []string) *PopupMenu {
	pm := &PopupMenu{}
	pm.InitPopupMenu(x, y, items)
	return pm
}

func (pm *PopupMenu) InitPopupMenu(x, y int, items []string) {
	pm.Items = items
	pm.ItemHeight = 24

	// サイズ計算（幅は固定、高さは項目数に応じて）
	width := 150
	height := len(items) * pm.ItemHeight

	pm.InitControlBase(pm, x, y, width, height)
	pm.InitDrawable(pm)
	pm.InitGrouping(pm)
	pm.ClippingFlag = false

	// 背景描画
	pm.OnDraw = func(screen *ebiten.Image) {
		theme := parts.CurrentTheme
		gx, gy := pm.GetGlobalPos()
		// 背景
		vector.FillRect(screen, float32(gx), float32(gy), float32(pm.Width), float32(pm.Height), theme.PopupBackground, false)
		// 枠線
		vector.StrokeRect(screen, float32(gx), float32(gy), float32(pm.Width), float32(pm.Height), 1, theme.FocusBorder, false)
	}

	// メニュー項目を生成
	for i, text := range items {
		item := newPopupMenuItem(pm, i, text, width, pm.ItemHeight)
		item.Y = i * pm.ItemHeight
		pm.AddChild(item)
	}
}

// Close はメニューを閉じる
func (pm *PopupMenu) Close() {
	PopupContainer.Close()
}

// --- PopupMenuItem ---

// popupMenuItemはPopupMenuの各項目
type popupMenuItem struct {
	InteractiveControl
	parts.TextDrawable

	menu  *PopupMenu
	index int
}

func newPopupMenuItem(menu *PopupMenu, index int, text string, width, height int) *popupMenuItem {
	theme := parts.CurrentTheme
	item := &popupMenuItem{}
	item.menu = menu
	item.index = index

	item.InitInteractiveControl(0, 0, width, height)
	item.InitTextDrawable(item, text, 14, parts.AlignLeft, parts.AlignCenter, 8, 0, theme.Text, true)

	// 描画
	item.OnDraw = func(screen *ebiten.Image) {
		gx, gy := item.GetGlobalPos()
		if item.IsMouseOver {
			vector.FillRect(screen, float32(gx), float32(gy), float32(item.Width), float32(item.Height), theme.PopupHover, false)
		}
	}

	// クリック
	item.OnRelease = func() {
		if menu.OnSelect != nil {
			menu.OnSelect(item.index)
		}
		menu.Close()
	}

	return item
}
