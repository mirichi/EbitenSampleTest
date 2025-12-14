package ui

import (
	"MyProject/ui/parts"
)

// TitleBarはウィンドウのタイトルバーを表すControl
// マウスドラッグによるウィンドウの移動機能を提供する
type TitleBar struct {
	InteractiveControl
	parts.TextDrawable
}

// NewTitleBarは新しいタイトルバーを作成する
// ドラッグ操作によるウィンドウ移動のイベントハンドラを設定する
func NewTitleBar(x, y, w, h int, text string) *TitleBar {
	t := &TitleBar{}
	t.InitTitleBar(x, y, w, h, text)
	return t
}

func (t *TitleBar) InitTitleBar(x, y, w, h int, text string) {
	theme := parts.CurrentTheme
	t.InitInteractiveControl(x, y, w, h)
	t.InitTextDrawable(t, text, h*2/3, parts.AlignCenter, parts.AlignCenter, 0, 0, theme.Text, true)
	t.BackColor = theme.TitleBar

	// TitleBarをドラッグすると親のWindowが移動するように設定
	var dragOffsetX, dragOffsetY int
	t.OnDragStart = func(x, y int) {
		// 親のWindowとの相対座標を保存
		cb := t.Parent.GetControlBase()
		dragOffsetX = x - cb.X
		dragOffsetY = y - cb.Y
	}
	t.OnDrag = func(x, y int) {
		// 親のWindowの座標を更新
		cb := t.Parent.GetControlBase()
		cb.X = x - dragOffsetX
		cb.Y = y - dragOffsetY
	}
}

// ClientAreaはウィンドウのクライアント領域を表すControl
// ウィンドウ内に配置される他のControlを保持する
type ClientArea struct {
	InteractiveControl
	parts.Grouping
}

// NewClientAreaは新しいクライアント領域を作成する
// ウィンドウのドラッグ移動に対応する
func NewClientArea(x, y, w, h int) *ClientArea {
	c := &ClientArea{}
	c.InitClientArea(x, y, w, h)
	return c
}

func (c *ClientArea) InitClientArea(x, y, w, h int) {
	c.InitInteractiveControl(x, y, w, h)
	c.InitGrouping(c)
	c.BackColor = parts.CurrentTheme.ClientArea
	c.AutoResizable = true

	// ClientAreaをドラッグすると親のWindowが移動する
	var dragOffsetX, dragOffsetY int
	c.OnDragStart = func(x, y int) {
		// 親のWindowとの相対座標を保存
		cb := c.Parent.GetControlBase()
		dragOffsetX = x - cb.X
		dragOffsetY = y - cb.Y
	}
	c.OnDrag = func(x, y int) {
		// 親のWindowの座標を更新
		cb := c.Parent.GetControlBase()
		cb.X = x - dragOffsetX
		cb.Y = y - dragOffsetY
	}
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
func NewWindow(x, y, w, h int, text string) *Window {
	win := &Window{}
	win.InitControlBase(win, x, y, w, h)
	win.InitResizable(win)
	win.InitGrouping(win)
	win.AutoLayout = parts.AutoLayoutFitV

	// TitleBarとClientArea生成
	win.TitleBar.InitTitleBar(0, 0, w, 30, text)
	win.ClientArea.InitClientArea(0, 30, w, h-30)

	// win.Groupingを指定することで、ウィンドウのレイアウトに追加される
	win.Grouping.AddChild(&win.TitleBar)
	win.Grouping.AddChild(&win.ClientArea)

	return win
}

// Windowに対してのAddChildはクライアント領域に委譲する
func (win *Window) AddChild(c parts.Control) {
	win.ClientArea.AddChild(c)
}
