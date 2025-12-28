package main

import (
	"math/rand"
	"time"
)

type CellState int

const (
	CellStateHidden CellState = iota
	CellStateOpened
	CellStateFlagged
)

type CellData struct {
	IsBomb        bool
	State         CellState
	NeighborCount int
}

type MinesweeperLogic struct {
	Width       int
	Height      int
	BombCount   int
	Grid        []CellData
	GameOver    bool
	GameWon     bool
	OpenedCount int
}

func NewMinesweeperLogic(w, h, bombs int) *MinesweeperLogic {
	l := &MinesweeperLogic{
		Width:     w,
		Height:    h,
		BombCount: bombs,
		Grid:      make([]CellData, w*h),
	}
	l.Reset()
	return l
}

func (l *MinesweeperLogic) index(x, y int) int {
	return y*l.Width + x
}

func (l *MinesweeperLogic) Reset() {
	for i := range l.Grid {
		l.Grid[i] = CellData{State: CellStateHidden}
	}

	l.GameOver = false
	l.GameWon = false
	l.OpenedCount = 0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	count := 0
	for count < l.BombCount {
		x := r.Intn(l.Width)
		y := r.Intn(l.Height)
		idx := l.index(x, y)
		if !l.Grid[idx].IsBomb {
			l.Grid[idx].IsBomb = true
			count++
		}
	}

	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			idx := l.index(x, y)
			if l.Grid[idx].IsBomb {
				continue
			}
			l.Grid[idx].NeighborCount = l.countBombs(x, y)
		}
	}
}

func (l *MinesweeperLogic) countBombs(cx, cy int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := cx+dx, cy+dy
			if nx >= 0 && nx < l.Width && ny >= 0 && ny < l.Height {
				if l.Grid[l.index(nx, ny)].IsBomb {
					count++
				}
			}
		}
	}
	return count
}

func (l *MinesweeperLogic) Open(x, y int) {
	idx := l.index(x, y)
	if l.GameOver || l.GameWon || l.Grid[idx].State != CellStateHidden {
		return
	}

	if l.Grid[idx].IsBomb {
		l.GameOver = true
		l.Grid[idx].State = CellStateOpened
		return
	}

	l.recursiveOpen(x, y)

	if l.OpenedCount == l.Width*l.Height-l.BombCount {
		l.GameWon = true
	}
}

func (l *MinesweeperLogic) recursiveOpen(x, y int) {
	if x < 0 || x >= l.Width || y < 0 || y >= l.Height {
		return
	}
	idx := l.index(x, y)
	if l.Grid[idx].State != CellStateHidden {
		return
	}

	l.Grid[idx].State = CellStateOpened
	l.OpenedCount++

	if l.Grid[idx].NeighborCount == 0 {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				l.recursiveOpen(x+dx, y+dy)
			}
		}
	}
}

func (l *MinesweeperLogic) ToggleFlag(x, y int) {
	if l.GameOver || l.GameWon {
		return
	}
	idx := l.index(x, y)
	if l.Grid[idx].State == CellStateHidden {
		l.Grid[idx].State = CellStateFlagged
	} else if l.Grid[idx].State == CellStateFlagged {
		l.Grid[idx].State = CellStateHidden
	}
}

func (l *MinesweeperLogic) GetFlagCount() int {
	count := 0
	for _, cell := range l.Grid {
		if cell.State == CellStateFlagged {
			count++
		}
	}
	return count
}

func (l *MinesweeperLogic) GetCell(x, y int) *CellData {
	if x < 0 || x >= l.Width || y < 0 || y >= l.Height {
		return nil
	}
	return &l.Grid[l.index(x, y)]
}
