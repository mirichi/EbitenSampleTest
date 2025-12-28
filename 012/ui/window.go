package ui

import (
	"MyProject/parts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TitleBarはウィンドウのタイトルバーを表すControl
// マウスドラッグによるウィンドウの移動機能を提供する
type TitleBar struct {
	InteractiveControl
	parts.TextDrawable
}

// NewTitleBarは新しいタイトルバーを作成する
// ドラッグ操作によるウィンドウ移動のイベントハンドラを設定する
func NewTitleBar(x, y, h float64, text string) *TitleBar {
	t := &TitleBar{}
	t.InitTitleBar(x, y, h, text)
	return t
}

func (t *TitleBar) InitTitleBar(x, y, h float64, text string) {
	theme := parts.CurrentTheme
	t.InitInteractiveControl(x, y, 0, h)
	t.InitTextDrawable(t, text, h*2/3, parts.AlignCenter, parts.AlignCenter, 0, 0, theme.Text, true)

	// TitleBarをドラッグすると親のWindowが移動するように設定
	var dragOffsetX, dragOffsetY float64
	t.OnDragStart = func(x, y float64) {
		// 親のWindowとの相対座標を保存
		cb := t.Parent.GetEntityBase()
		dragOffsetX = x - cb.X
		dragOffsetY = y - cb.Y
	}
	t.OnDrag = func(x, y float64) {
		// 親のWindowの座標を更新
		cb := t.Parent.GetEntityBase()
		cb.X = x - dragOffsetX
		cb.Y = y - dragOffsetY
	}

	t.AddBeforeDrawFunction(func(screen *ebiten.Image) {
		gx, gy := t.GetGlobalPos()
		vector.FillRect(screen, float32(gx), float32(gy), float32(t.Width), float32(t.Height), theme.TitleBar, false)
	})
}

// ClientAreaはウィンドウのクライアント領域を表すControl
// ウィンドウ内に配置される他のControlを保持する
type ClientArea struct {
	InteractiveControl
	parts.Grouping
}

// NewClientAreaは新しいクライアント領域を作成する
// ウィンドウのドラッグ移動に対応する
func NewClientArea(x, y float64) *ClientArea {
	c := &ClientArea{}
	c.InitClientArea(x, y)
	return c
}

func (c *ClientArea) InitClientArea(x, y float64) {
	c.InitInteractiveControl(x, y, 0, 0)
	c.InitGrouping(c)
	c.FlexGrow = 1

	// ClientAreaをドラッグすると親のWindowが移動する
	var dragOffsetX, dragOffsetY float64
	c.OnDragStart = func(x, y float64) {
		// 親のWindowとの相対座標を保存
		cb := c.Parent.GetEntityBase()
		dragOffsetX = x - cb.X
		dragOffsetY = y - cb.Y
	}
	c.OnDrag = func(x, y float64) {
		// 親のWindowの座標を更新
		cb := c.Parent.GetEntityBase()
		cb.X = x - dragOffsetX
		cb.Y = y - dragOffsetY
	}

	c.AddBeforeDrawFunction(func(screen *ebiten.Image) {
		gx, gy := c.GetGlobalPos()
		vector.FillRect(screen, float32(gx), float32(gy), float32(c.Width), float32(c.Height), parts.CurrentTheme.ClientArea, false)
	})
}

// Windowはタイトルバーとクライアント領域を持つ複合Control
// これ自体はコンテナとしての役割を持ち、具体的な機能はTitleBarとClientAreaに実装する
type Window struct {
	parts.ControlBase
	parts.Resizable // Resizableは自前でマウス処理を持っている
	parts.Grouping

	TitleBar   TitleBar
	ClientArea ClientArea
}

// NewWindowは新しいウィンドウを作成する
// タイトルバーとクライアント領域を初期化し、レイアウトを設定する
func NewWindow(x, y, w, h float64, text string) *Window {
	win := &Window{}
	win.InitControlBase(win, x, y, w, h)
	win.InitResizable(win)
	win.InitGrouping(win)
	win.AutoLayout = parts.FlexLayoutV(parts.FlexStart, parts.FlexStretch, 0)
	win.ClippingFlag = true

	// TitleBarとClientArea生成
	win.TitleBar.InitTitleBar(0, 0, 30, text)
	win.ClientArea.InitClientArea(0, 30)

	// win.Groupingを指定することで、ウィンドウのレイアウトに追加される
	win.Grouping.AddChild(&win.TitleBar)
	win.Grouping.AddChild(&win.ClientArea)

	return win
}

// Windowに対してのAddChildはクライアント領域に委譲する
func (win *Window) AddChild(c parts.Entity) {
	win.ClientArea.AddChild(c)
}
