package parts

import (
	"MyProject/ui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

// ControlBaseを埋め込むとControlインターフェースを満たす
type Control interface {
	HandleInput(t input.Touch) bool
	Update()
	Draw(screen *ebiten.Image)
	GetControlBase() *ControlBase
}

// ControlBaseはコントロールの基本機能を実装した構造体
// 親コントロールの参照、座標、サイズ、可視性などの共通プロパティを管理する
// また、HandleInput/Update/Drawの各フェーズで実行される関数リストを保持している
type ControlBase struct {
	Control              Control
	Parent               Control
	X, Y                 int
	Width, Height        int
	Visible              bool
	AutoResizable        bool
	handleInputFunctions []func(t input.Touch) bool
	updateFunctions      []func()
	drawFunctions        []func(screen *ebiten.Image)
}

func NewControlBase(c Control, x, y, w, h int) *ControlBase {
	cb := &ControlBase{Control: c, X: x, Y: y, Width: w, Height: h, Visible: true, AutoResizable: false}
	return cb
}

// AddHandleInputFunctionは、HandleInputメソッド呼び出し時に実行される関数を登録する
// 登録された関数は後に追加されたものから順に実行される
func (c *ControlBase) AddHandleInputFunction(f func(t input.Touch) bool) {
	c.handleInputFunctions = append(c.handleInputFunctions, f)
}

// 登録されたHandleInputFunctionsを順次実行する
func (c *ControlBase) HandleInput(t input.Touch) bool {
	if c.Visible && t != nil {
		for i := len(c.handleInputFunctions) - 1; i >= 0; i-- {
			if c.handleInputFunctions[i](t) {
				return true
			}
		}
	}

	return false
}

// AddUpdateFunctionは、Updateメソッド呼び出し時に実行される関数を登録する
// 登録された関数は後に追加されたものから順に実行される
func (c *ControlBase) AddUpdateFunction(f func()) {
	c.updateFunctions = append(c.updateFunctions, f)
}

// 登録されたUpdateFunctionsを順次実行する
func (c *ControlBase) Update() {
	if c.Visible {
		for i := len(c.updateFunctions) - 1; i >= 0; i-- {
			c.updateFunctions[i]()
		}
	}
}

// AddDrawFunctionは、Drawメソッド呼び出し時に実行される関数を登録する
// 登録された関数は登録順に実行される
func (c *ControlBase) AddDrawFunction(f func(screen *ebiten.Image)) {
	c.drawFunctions = append(c.drawFunctions, f)
}

// 登録されたDrawFunctionsを順次実行する
func (c *ControlBase) Draw(screen *ebiten.Image) {
	if c.Visible {
		for _, f := range c.drawFunctions {
			f(screen)
		}
	}
}

// getGlobalPosは、画面全体（ルートコントロール）からの絶対座標を取得
// 親コントロールの座標を再帰的に加算して算出する
func (c *ControlBase) GetGlobalPos() (int, int) {
	p := c.Parent
	if p == nil {
		return c.X, c.Y
	} else {
		pc := p.GetControlBase()
		x, y := pc.GetGlobalPos()
		return x + c.X, y + c.Y
	}
}

// ControlBaseのポインタを返す。Controlインターフェースのメソッドの一つ
// Controlとしての基本機能はControlBaseに定義することでControlインターフェースをなるべくシンプルにしておきたい
func (c *ControlBase) GetControlBase() *ControlBase {
	return c
}
