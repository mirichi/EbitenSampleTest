package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TitleBarはウィンドウのタイトルバーのコントロール
// ドラッグできる
type TitleBar struct {
	*ControlBase
	*Drawable
	*TextDrawable
	*MouseInteraction
}

// ClientAreaはウィンドウのクライアント領域のコントロール
// ウィンドウ上のコントロールを保持する
type ClientArea struct {
	*ControlBase
	*Drawable
	*MouseInteraction
	*Grouping
}

// Windowはタイトルバーとクライアント領域を持つウィンドウコントロール
// これ自体には機能は無い
type Window struct {
	*ControlBase
	*Resizable
	*Grouping

	*TitleBar
	*ClientArea
}

// TitleBar生成
func newTitleBar(x, y, w, h int, text string) *TitleBar {
	t := &TitleBar{}
	t.ControlBase = NewControlBase(t, x, y, w, h)
	t.Drawable = NewDrawable(t, t.drawTitleBar)
	t.TextDrawable = NewTextDrawable(t, text, h*2/3, AlignCenter, AlignCenter, 0, 0, color.White, true)

	// TitleBarをドラッグすると親のWindowが移動する
	t.MouseInteraction = NewMouseInteraction(t)
	t.OnDrag = func(dx, dy int) {
		t.Parent.GetControlBase().X += dx
		t.Parent.GetControlBase().Y += dy
	}

	return t
}

func (t *TitleBar) drawTitleBar(screen *ebiten.Image) {
	ox, oy := t.GetOwnerPos()
	vector.FillRect(screen, float32(ox+t.X), float32(oy+t.Y), float32(t.Width), float32(t.Height), color.RGBA{0x00, 0x40, 0x00, 0xff}, false)
}

// ClientArea生成
func newClientArea(x, y, w, h int) *ClientArea {
	c := &ClientArea{}
	c.ControlBase = NewControlBase(c, x, y, w, h)
	c.AutoResizable = true

	// DrawFunctionsは前から順に実行
	// Groupingより後に実行するとボタンがClientAreaに上書きされて消える
	c.Drawable = NewDrawable(c, c.drawClientArea)

	// ClientAreaをドラッグすると親のWindowが移動する
	c.MouseInteraction = NewMouseInteraction(c)
	c.OnDrag = func(dx, dy int) {
		c.Parent.GetControlBase().X += dx
		c.Parent.GetControlBase().Y += dy
	}

	// UpdateFunctionsは後ろから順に実行
	// Draggableより後に実行すると操作がClientAreaに食われてボタンを押せない
	c.Grouping = NewGrouping(c.ControlBase)

	return c
}

func (c *ClientArea) drawClientArea(screen *ebiten.Image) {
	ox, oy := c.GetOwnerPos()
	vector.FillRect(screen, float32(ox+c.X), float32(oy+c.Y), float32(c.Width), float32(c.Height), color.RGBA{0x30, 0x30, 0x30, 0xff}, false)
}

// Window生成
func NewWindow(x, y, w, h int, text string) *Window {
	win := &Window{}
	win.ControlBase = NewControlBase(win, x, y, w, h)
	win.Grouping = NewGrouping(win.ControlBase)
	win.Resizable = NewResizable(win.ControlBase)
	win.OnResize = func() {
		win.Layout()
	}
	win.AutoLayout = NewAutoLayoutFitV(win.Grouping)

	// TitleBarとClientArea生成
	win.TitleBar = newTitleBar(0, 0, w, 30, text)
	win.ClientArea = newClientArea(0, 30, w, h-30)

	// Groupingを指定しないと↓のメソッドが呼ばれてハマる(ハマった)
	win.Grouping.AddChild(win.TitleBar)
	win.Grouping.AddChild(win.ClientArea)

	return win
}

// Windowに対してのAddChildはクライアント領域に委譲する
func (win *Window) AddChild(c Control) {
	win.ClientArea.AddChild(c)
}
