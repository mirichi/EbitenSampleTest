package ui

import (
	"MyProject/ui/parts"
	"image/color"
)

// TitleBarはウィンドウのタイトルバーを表すWidget
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
	t.InitInteractiveControl(x, y, w, h)
	t.InitTextDrawable(t, text, h*2/3, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
	t.BackColor = color.RGBA{0x00, 0x40, 0x00, 0xff}

	// TitleBarをドラッグすると親のWindowが移動するように設定
	var dragOffsetX, dragOffsetY int
	t.OnDragStart = func(x, y int) {
		// ドラッグ開始時のマウス位置とウィンドウ位置のオフセットを保存
		cb := t.Parent.GetControlBase()
		gx, gy := cb.GetGlobalPos()
		dragOffsetX = x - gx
		dragOffsetY = y - gy
	}
	t.OnDrag = func(x, y int) {
		// マウス位置からオフセットを引いてウィンドウの新しい位置を計算
		// さらに親コントロール（通常はMainScreen）の座標を引いてローカル座標に変換
		cb := t.Parent.GetControlBase()
		ox, oy := cb.Parent.GetControlBase().GetGlobalPos()
		cb.X = x - dragOffsetX - ox
		cb.Y = y - dragOffsetY - oy
	}
}

// ClientAreaはウィンドウのクライアント領域を表すWidget
// ウィンドウ内に配置される他のWidgetを保持する
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
	c.BackColor = color.RGBA{0x30, 0x30, 0x30, 0xff}
	c.AutoResizable = true

	// ClientAreaをドラッグすると親のWindowが移動する
	var dragOffsetX, dragOffsetY int
	c.OnDragStart = func(x, y int) {
		// 親のWindowの座標を保存
		cb := c.Parent.GetControlBase()
		gx, gy := cb.GetGlobalPos()
		dragOffsetX = x - gx
		dragOffsetY = y - gy
	}
	c.OnDrag = func(x, y int) {
		// 親のWindowの座標を更新
		cb := c.Parent.GetControlBase()
		ox, oy := cb.Parent.GetControlBase().GetGlobalPos()
		cb.X = x - dragOffsetX - ox
		cb.Y = y - dragOffsetY - oy
	}
}

// Windowはタイトルバーとクライアント領域を持つ複合Widget
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
	win.InitGrouping(win)
	win.InitResizable(win)
	win.AutoLayout = parts.AutoLayoutFitV

	// TitleBarとClientArea生成
	win.TitleBar.InitTitleBar(0, 0, w, 30, text)
	win.ClientArea.InitClientArea(0, 30, w, h-30)

	// win.Groupingを指定することで、ウィンドウのレイアウトに追加される
	win.Grouping.AddChild(&win.TitleBar)
	win.Grouping.AddChild(&win.ClientArea)

	win.OnResize = func() {
		win.Layout()
	}

	return win
}

// Windowに対してのAddChildはクライアント領域に委譲する
func (win *Window) AddChild(c parts.Control) {
	win.ClientArea.AddChild(c)
}
