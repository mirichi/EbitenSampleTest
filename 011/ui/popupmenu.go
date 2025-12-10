package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PopupMenuはドロップダウンやコンテキストメニューに使用するポップアップメニュー
type PopupMenu struct {
	parts.ControlBase
	parts.Drawable
	parts.Grouping

	Items       []string
	OnSelect    func(index int)
	OnClose     func()
	ItemHeight  int
	backColor   color.Color
	hoverColor  color.Color
	borderColor color.Color
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
	pm.backColor = color.RGBA{0x40, 0x40, 0x40, 0xff}
	pm.hoverColor = color.RGBA{0x60, 0x80, 0xc0, 0xff}
	pm.borderColor = color.RGBA{0x80, 0x80, 0x80, 0xff}

	// サイズ計算（幅は固定、高さは項目数に応じて）
	width := 150
	height := len(items) * pm.ItemHeight

	pm.InitControlBase(pm, x, y, width, height)
	pm.InitDrawable(pm)
	pm.InitGrouping(pm)
	pm.ClippingFlag = false

	// 背景描画
	pm.OnDraw = func(screen *ebiten.Image) {
		gx, gy := pm.GetGlobalPos()
		// 背景
		vector.FillRect(screen, float32(gx), float32(gy), float32(pm.Width), float32(pm.Height), pm.backColor, false)
		// 枠線
		vector.StrokeRect(screen, float32(gx), float32(gy), float32(pm.Width), float32(pm.Height), 1, pm.borderColor, false)
	}

	// メニュー項目を生成
	for i, text := range items {
		item := newPopupMenuItem(pm, i, text, width, pm.ItemHeight)
		item.Y = i * pm.ItemHeight
		pm.AddChild(item)
	}

	// // ホバー状態を解除
	// pm.OnBeforeHandleInput = func() {
	// 	for _, child := range pm.Children {
	// 		child.(*popupMenuItem).IsHovering = false
	// 	}
	// }
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

	menu       *PopupMenu
	index      int
	IsHovering bool
}

func newPopupMenuItem(menu *PopupMenu, index int, text string, width, height int) *popupMenuItem {
	item := &popupMenuItem{}
	item.menu = menu
	item.index = index

	item.InitInteractiveControl(0, 0, width, height)
	item.InitTextDrawable(item, text, 14, parts.AlignLeft, parts.AlignCenter, 8, 0, color.White, true)

	// 描画
	item.OnDraw = func(screen *ebiten.Image) {
		gx, gy := item.GetGlobalPos()
		if item.IsHovering {
			vector.FillRect(screen, float32(gx), float32(gy), float32(item.Width), float32(item.Height), menu.hoverColor, false)
		}
	}

	// ホバー状態を解除
	item.OnBeforeHandleInput = func() {
		item.IsHovering = false
	}

	// ホバー
	item.OnHover = func() {
		item.IsHovering = true
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
