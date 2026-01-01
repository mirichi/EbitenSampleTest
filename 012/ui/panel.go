package ui

import (
	"image/color"

	"MyProject/parts"
	"MyProject/uiparts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Panel は角丸の背景を持つシンプルなコンテナControl
type Panel struct {
	uiparts.ControlBase
	parts.Grouping
	BackgroundColor color.Color
	CornerRadius    float64
}

// NewPanel は新しいパネルを作成する
func NewPanel(x, y, w, h float64) *Panel {
	p := &Panel{}
	p.InitPanel(x, y, w, h)
	return p
}

func (p *Panel) InitPanel(x, y, w, h float64) {
	p.InitControlBase(p, x, y, w, h)
	p.InitGrouping(p)
	p.BackgroundColor = color.RGBA{0, 0, 0, 255} // 黒背景
	p.CornerRadius = 10

	p.AddBeforeDrawFunction(p.drawBackground)
}

func (p *Panel) drawBackground(screen *ebiten.Image) {
	gx, gy := p.GetGlobalPos()
	x, y := float32(gx), float32(gy)
	w, h := float32(p.Width), float32(p.Height)
	r := float32(p.CornerRadius)
	c := p.BackgroundColor

	// 角丸の描画
	// 左上
	vector.DrawFilledCircle(screen, x+r, y+r, r, c, true)
	// 右上
	vector.DrawFilledCircle(screen, x+w-r, y+r, r, c, true)
	// 左下
	vector.DrawFilledCircle(screen, x+r, y+h-r, r, c, true)
	// 右下
	vector.DrawFilledCircle(screen, x+w-r, y+h-r, r, c, true)

	// 中央の矩形（縦長） - 上辺から下辺まで、角の分を除いた幅
	vector.DrawFilledRect(screen, x+r, y, w-2*r, h, c, true)
	// 左右の矩形 - 角の分を除いた高さ
	vector.DrawFilledRect(screen, x, y+r, r, h-2*r, c, true)     // 左側
	vector.DrawFilledRect(screen, x+w-r, y+r, r, h-2*r, c, true) // 右側
}
