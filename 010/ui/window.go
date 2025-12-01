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
	*parts.ControlBase
	*parts.Drawable
	*parts.TextDrawable
	*parts.MouseInteraction
}

// ClientAreaはウィンドウのクライアント領域を表すコントロール
// ウィンドウ内に配置される他のコントロールを保持する
type ClientArea struct {
	*parts.ControlBase
	*parts.Drawable
	*parts.MouseInteraction
	*parts.Grouping
}

// Windowはタイトルバーとクライアント領域を持つウィンドウコントロール
// これ自体はコンテナとしての役割を持ち、具体的な機能はTitleBarとClientAreaに委譲する
type Window struct {
	*parts.ControlBase
	*parts.Resizable // Resizableは自前でマウス処理を持っている
	*parts.Grouping

	TitleBar   *TitleBar
	ClientArea *ClientArea
}

// newTitleBarは新しいタイトルバーを作成する
// ドラッグ操作によるウィンドウ移動のイベントハンドラを設定する
func newTitleBar(x, y, w, h int, text string) *TitleBar {
	t := &TitleBar{}
	t.ControlBase = parts.NewControlBase(t, x, y, w, h)
	t.Drawable = parts.NewDrawable(t, t.drawTitleBar)
	t.TextDrawable = parts.NewTextDrawable(t, text, h*2/3, parts.AlignCenter, parts.AlignCenter, 0, 0, color.White, true)

	// TitleBarをドラッグすると親のWindowが移動するように設定
	t.MouseInteraction = parts.NewMouseInteraction(t)
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

	return t
}

func (t *TitleBar) drawTitleBar(screen *ebiten.Image) {
	gx, gy := t.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(t.Width), float32(t.Height), color.RGBA{0x00, 0x40, 0x00, 0xff}, false)
}

// newClientAreaは新しいクライアント領域を作成する
// ウィンドウのリサイズやドラッグ移動に対応する
func newClientArea(x, y, w, h int) *ClientArea {
	c := &ClientArea{}
	c.ControlBase = parts.NewControlBase(c, x, y, w, h)
	c.AutoResizable = true

	// DrawFunctionsは前から順に実行
	// Groupingより後に実行するとボタンがClientAreaに上書きされて消える
	c.Drawable = parts.NewDrawable(c, c.drawClientArea)

	// ClientAreaをドラッグすると親のWindowが移動する
	c.MouseInteraction = parts.NewMouseInteraction(c)
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

	// UpdateFunctionsは後ろから順に実行
	// Draggableより後に実行すると操作がClientAreaに食われてボタンを押せない
	c.Grouping = parts.NewGrouping(c.ControlBase)

	return c
}

func (c *ClientArea) drawClientArea(screen *ebiten.Image) {
	gx, gy := c.GetGlobalPos()
	vector.FillRect(screen, float32(gx), float32(gy), float32(c.Width), float32(c.Height), color.RGBA{0x30, 0x30, 0x30, 0xff}, false)
}

// NewWindowは新しいウィンドウを作成する
// タイトルバーとクライアント領域を初期化し、レイアウトを設定する
func NewWindow(x, y, w, h int, text string) *Window {
	win := &Window{}
	win.ControlBase = parts.NewControlBase(win, x, y, w, h)
	win.Grouping = parts.NewGrouping(win.ControlBase)
	win.Resizable = parts.NewResizable(win.ControlBase)
	win.OnResize = func() {
		win.Layout()
	}
	win.AutoLayout = parts.NewAutoLayoutFitV(win.Grouping)

	// TitleBarとClientArea生成
	win.TitleBar = newTitleBar(0, 0, w, 30, text)
	win.ClientArea = newClientArea(0, 30, w, h-30)

	// ClipGroupingを指定しないと↓のメソッドが呼ばれてハマる(ハマった)
	win.Grouping.AddChild(win.TitleBar)
	win.Grouping.AddChild(win.ClientArea)

	return win
}

// Windowに対してのAddChildはクライアント領域に委譲する
func (win *Window) AddChild(c parts.Control) {
	win.ClientArea.AddChild(c)
}
