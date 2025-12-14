package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"MyProject/ui/parts"
)

// Canvasは描画可能なキャンバス
type Canvas struct {
	parts.ControlBase
	parts.ImageDrawable
	parts.MouseInteraction

	lastX     int
	lastY     int
	Color     color.Color
	LineWidth int
}

// Canvas生成
func NewCanvas(x, y, w, h int) *Canvas {
	c := &Canvas{}
	c.InitControlBase(c, x, y, w, h)
	c.InitImageDrawable(c, ebiten.NewImage(w, h))
	c.InitMouseInteraction(c)
	c.Color = color.White
	c.LineWidth = 6

	c.Clear()

	c.OnDragStart = func(x, y int) {
		// キャンバス内の相対座標に変換
		gx, gy := c.GetGlobalPos()
		lx, ly := x-gx, y-gy
		c.lastX, c.lastY = lx, ly

		// 点を描画
		vector.FillCircle(c.Image, float32(lx), float32(ly), float32(c.LineWidth/2), c.Color, true)
	}

	c.OnDrag = func(x, y int) {
		// キャンバス内の相対座標に変換
		gx, gy := c.GetGlobalPos()
		cx, cy := x-gx, y-gy

		// 線を描画
		vector.StrokeLine(c.Image, float32(cx), float32(cy), float32(c.lastX), float32(c.lastY), float32(c.LineWidth), c.Color, true)
		vector.FillCircle(c.Image, float32(c.lastX), float32(c.lastY), float32(c.LineWidth/2), c.Color, true)

		c.lastX, c.lastY = cx, cy
	}

	return c
}

func (c *Canvas) Clear() {
	c.Image.Fill(color.Black)
}
