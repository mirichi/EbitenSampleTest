package ui

import "MyProject/uiparts"

// BlankControlは位置とサイズだけを持つ空のコントロール
type BlankControl struct {
	uiparts.ControlBase
}

func NewBlankControl(x, y, w, h float64) *BlankControl {
	b := &BlankControl{}
	b.InitBlankControl(x, y, w, h)
	return b
}

func (b *BlankControl) InitBlankControl(x, y, w, h float64) {
	b.InitControlBase(b, x, y, w, h)
}
