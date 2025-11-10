package ui

import "github.com/hajimehoshi/ebiten/v2"

// ControlBaseを埋め込むとControlインターフェースを満たす
type Control interface {
	Update(t TouchInfo)
	Draw(screen *ebiten.Image)
	GetControlBase() *ControlBase
}

// ControlBaseはコントロールの基本機能を実装した構造体
// ・親となるコントロールを持つ
// ・座標とサイズを持つ(使わないこともあるがGUI用コントロールなので)
// ・Update/Drawのメソッドと、メソッド呼び出し時に実行する関数のスライスを持つ
type ControlBase struct {
	Parent          *ControlBase
	X, Y            int
	Width, Height   int
	updateFunctions []func(t TouchInfo)
	drawFunctions   []func(screen *ebiten.Image)
}

// ControlBase生成
func NewControlBase(x, y, w, h int) *ControlBase {
	return &ControlBase{X: x, Y: y, Width: w, Height: h}
}

// このコントロールのUpdateを呼び出した際に実行するUpdateFunctionを登録する
func (c *ControlBase) AddUpdateFunction(f func(t TouchInfo)) {
	c.updateFunctions = append(c.updateFunctions, f)
}

// 登録されたUpdateFunctionsを順次実行する
func (c *ControlBase) Update(t TouchInfo) {
	for _, f := range c.updateFunctions {
		f(t)
	}
}

// このコントロールのDrawを呼び出した際に実行するDrawFunctionを登録する
func (c *ControlBase) AddDrawFunction(f func(screen *ebiten.Image)) {
	c.drawFunctions = append(c.drawFunctions, f)
}

// 登録されたDrawFunctionsを順次実行する
func (c *ControlBase) Draw(screen *ebiten.Image) {
	for _, f := range c.drawFunctions {
		f(screen)
	}
}

// 親Controlの座標を取得する
func (c *ControlBase) GetOwnerPos() (int, int) {
	p := c.Parent
	if p == nil {
		return 0, 0
	} else {
		x, y := p.GetOwnerPos()
		return x + p.X, y + p.Y
	}
}

// ControlBaseのポインタを返す。Controlインターフェースのメソッドの一つ。
// Controlとしての基本機能はControlBaseに定義することでControlインターフェースをなるべくシンプルにしておきたい。
func (c *ControlBase) GetControlBase() *ControlBase {
	return c
}
