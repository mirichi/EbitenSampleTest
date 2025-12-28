package main

import (
	"image/color"
	"strconv"

	"MyProject/parts"
	"MyProject/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Cell struct {
	ui.InteractiveControl
	parts.TextDrawable
	parts.Focusable

	XIndex    int
	YIndex    int
	Logic     *MinesweeperLogic
	BackColor color.Color
}

func NewCell(x, y, w, h float64, xi, yi int, logic *MinesweeperLogic) *Cell {
	c := &Cell{
		XIndex: xi,
		YIndex: yi,
		Logic:  logic,
	}
	c.InitCell(x, y, w, h)
	return c
}

func (c *Cell) InitCell(x, y, w, h float64) {
	theme := parts.CurrentTheme

	c.InitInteractiveControl(x, y, w, h)
	c.InitTextDrawable(c, "", 16, parts.AlignCenter, parts.AlignCenter, 0, 0, theme.Text, false)
	c.InitFocusable(c)

	c.OnClick = func() {
		c.Logic.Open(c.XIndex, c.YIndex)
	}

	c.OnRightRelease = func() {
		c.Logic.ToggleFlag(c.XIndex, c.YIndex)
	}

	c.AddAfterUpdateFunction(func() {
		cellData := c.Logic.GetCell(c.XIndex, c.YIndex)
		if cellData == nil {
			return
		}
		theme := parts.CurrentTheme

		switch cellData.State {
		case CellStateHidden:
			c.Text = ""
			if c.IsHovering {
				c.BackColor = theme.ButtonHover
			} else {
				c.BackColor = theme.ButtonNormal
			}
		case CellStateFlagged:
			c.Text = "P"
			c.TextColor = color.RGBA{0xff, 0x00, 0x00, 0xff}
			c.BackColor = theme.ButtonNormal
		case CellStateOpened:
			c.BackColor = color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
			if cellData.IsBomb {
				c.Text = "*"
				c.TextColor = color.Black
				c.BackColor = color.RGBA{0xff, 0x88, 0x88, 0xff}
			} else {
				if cellData.NeighborCount > 0 {
					c.Text = strconv.Itoa(cellData.NeighborCount)
					c.TextColor = getNumberColor(cellData.NeighborCount)
				} else {
					c.Text = ""
				}
			}
		}
	})

	c.AddBeforeDrawFunction(func(screen *ebiten.Image) {
		gx, gy := c.GetGlobalPos()
		vector.FillRect(screen, float32(gx), float32(gy), float32(c.Width), float32(c.Height), c.BackColor, false)
	})
}

func getNumberColor(n int) color.Color {
	colors := []color.RGBA{
		{0x00, 0x00, 0xff, 0xff}, // 1: Blue
		{0x00, 0x81, 0x00, 0xff}, // 2: Green
		{0xff, 0x00, 0x00, 0xff}, // 3: Red
		{0x00, 0x00, 0x81, 0xff}, // 4: Dark Blue
		{0x81, 0x00, 0x00, 0xff}, // 5: Maroon
		{0x00, 0x81, 0x81, 0xff}, // 6: Cyan
		{0x00, 0x00, 0x00, 0xff}, // 7: Black
		{0x81, 0x81, 0x81, 0xff}, // 8: Gray
	}
	if n > 0 && n <= len(colors) {
		return colors[n-1]
	}
	return color.Black
}
