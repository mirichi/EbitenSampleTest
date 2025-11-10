package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// TitleBarはウィンドウのタイトルバーのコントロール
// ドラッグできる
type TitleBar struct {
	*ControlBase
	*Draggable
	*Drawable
	*TextDrawable
}

// ClientAreaはウィンドウのクライアント領域のコントロール
// ウィンドウ上のコントロールを保持する
type ClientArea struct {
	*ControlBase
	*Drawable
	*Grouping
}

// Windowはタイトルバーとクライアント領域を持つウィンドウコントロール
// これ自体には機能は無い
type Window struct {
	*ControlBase
	*Grouping

	*TitleBar
	*ClientArea
}

// TitleBar生成
func newTitleBar(x, y, w, h int, text string) *TitleBar {
	image := ebiten.NewImage(w, h)
	image.Fill(color.RGBA{0x00, 0x40, 0x00, 0xff})

	cb := NewControlBase(x, y, w, h)

	// TitleBarをドラッグすると親のWindowが移動する
	dg := NewDraggable(cb,
		func(dx, dy int) {
			cb.Parent.X += dx
			cb.Parent.Y += dy
		},
	)

	t := &TitleBar{
		ControlBase:  cb,
		Draggable:    dg,
		Drawable:     NewDrawable(cb, image),
		TextDrawable: NewTextDrawable(cb, text, h*2/3, AlignCenter, AlignCenter, 0, 0, color.White, true),
	}
	return t
}

// ClientArea生成
func newClientArea(x, y, w, h int) *ClientArea {
	image := ebiten.NewImage(w, h)
	image.Fill(color.RGBA{0x30, 0x30, 0x30, 0xff})

	cb := NewControlBase(x, y, w, h)

	c := &ClientArea{
		ControlBase: cb,
		Drawable:    NewDrawable(cb, image),
		Grouping:    NewGrouping(cb),
	}
	return c
}

// Window生成
func NewWindow(x, y, w, h int, text string) *Window {
	cb := NewControlBase(x, y, w, h)
	win := &Window{
		ControlBase: cb,
		Grouping:    NewGrouping(cb),
	}

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
