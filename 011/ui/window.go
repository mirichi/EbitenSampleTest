package ui

import (
	"MyProject/ui/parts"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TitleBarはウィンドウのタイトルバーを表すコントロール
// マウスドラッグによるウィンドウの移動機能を提供する
type TitleBar struct {
	InteractiveWidget
	parts.TextDrawable
}

// NewTitleBarは新しいタイトルバーを作成する
// ドラッグ操作によるウィンドウ移動のイベントハンドラを設定する
func NewTitleBar(x, y, w, h int, text string) *TitleBar {
	t := &TitleBar{}
	t.InitTitleBar(nil, x, y, w, h, text)
	return t
}

func (t *TitleBar) InitTitleBar(g parts.AddChilder, x, y, w, h int, text string) {
	t.InitInteractiveWidget(nil, x, y, w, h)
	t.InitTextDrawable(t, text, h*2/3, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)
	if g != nil {
		g.AddChild(t)
	}
	t.OnDraw = t.drawTitleBar

	// TitleBarをドラッグすると親のWindowが移動するように設定
	var dragOffsetX, dragOffsetY int
	t.OnDragStart = func(x, y int) {
		// ドラッグ開始時のマウス位置とウィンドウ位置のオフセットを保存
		cb := t.Parent.GetWidgetBase()
		gx, gy := cb.GetGlobalPos()
		dragOffsetX = x - gx
		dragOffsetY = y - gy
	}
	t.OnDrag = func(x, y int) {
		// マウス位置からオフセットを引いてウィンドウの新しい位置を計算
		// さらに親コントロール（通常はMainScreen）の座標を引いてローカル座標に変換
		cb := t.Parent.GetWidgetBase()
		ox, oy := cb.Parent.GetWidgetBase().GetGlobalPos()
		cb.X = x - dragOffsetX - ox
		cb.Y = y - dragOffsetY - oy
	}
}

func (t *TitleBar) drawTitleBar(screen *ebiten.Image) {
	gx, gy := t.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(t.Width), float32(t.Height), color.RGBA{0x00, 0x40, 0x00, 0xff}, false)
}

// ClientAreaはウィンドウのクライアント領域を表すコントロール
// ウィンドウ内に配置される他のコントロールを保持する
type ClientArea struct {
	InteractiveWidget
	parts.Grouping
}

// NewClientAreaは新しいクライアント領域を作成する
// ウィンドウのドラッグ移動に対応する
func NewClientArea(x, y, w, h int) *ClientArea {
	c := &ClientArea{}
	c.InitClientArea(nil, x, y, w, h)
	return c
}

func (c *ClientArea) InitClientArea(g parts.AddChilder, x, y, w, h int) {
	c.InitInteractiveWidget(nil, x, y, w, h)
	c.InitGrouping(c)
	if g != nil {
		g.AddChild(c)
	}
	c.OnDraw = c.drawClientArea
	c.AutoResizable = true

	// ClientAreaをドラッグすると親のWindowが移動する
	var dragOffsetX, dragOffsetY int
	c.OnDragStart = func(x, y int) {
		// 親のWindowの座標を保存
		cb := c.Parent.GetWidgetBase()
		gx, gy := cb.GetGlobalPos()
		dragOffsetX = x - gx
		dragOffsetY = y - gy
	}
	c.OnDrag = func(x, y int) {
		// 親のWindowの座標を更新
		cb := c.Parent.GetWidgetBase()
		ox, oy := cb.Parent.GetWidgetBase().GetGlobalPos()
		cb.X = x - dragOffsetX - ox
		cb.Y = y - dragOffsetY - oy
	}
}

func (c *ClientArea) drawClientArea(screen *ebiten.Image) {
	gx, gy := c.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(c.Width), float32(c.Height), color.RGBA{0x30, 0x30, 0x30, 0xff}, false)
}

// Windowはタイトルバーとクライアント領域を持つ複合コントロール
// これ自体はコンテナとしての役割を持ち、具体的な機能はTitleBarとClientAreaに実装する
type Window struct {
	parts.WidgetBase
	parts.Resizable // Resizableは自前でマウス処理を持っている
	parts.Grouping

	TitleBar   TitleBar
	ClientArea ClientArea
}

// NewWindowは新しいウィンドウを作成する
// タイトルバーとクライアント領域を初期化し、レイアウトを設定する
func NewWindow(x, y, w, h int, text string) *Window {
	win := &Window{}
	win.InitWidgetBase(win, x, y, w, h)
	win.InitGrouping(win)
	win.InitResizable(win)
	win.AutoLayout = parts.AutoLayoutFitV

	// TitleBarとClientArea生成
	// win.Groupingを指定することで、ウィンドウのレイアウトに追加される
	win.TitleBar.InitTitleBar(&win.Grouping, 0, 0, w, 30, text)
	win.ClientArea.InitClientArea(&win.Grouping, 0, 30, w, h-30)

	win.OnResize = func() {
		win.Layout()
	}

	return win
}

// Windowに対してのAddChildはクライアント領域に委譲する
func (win *Window) AddChild(c parts.Widget) {
	win.ClientArea.AddChild(c)
}
