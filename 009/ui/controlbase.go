package ui

import "github.com/hajimehoshi/ebiten/v2"

// UpdateResultはUpdateメソッドの戻り値の型
type UpdateResult int

const (
	NotConsumed     UpdateResult = 0 // イベントは消費されず、親や他のコントロールに伝播する
	Consumed        UpdateResult = 1 // イベントは消費され、親への伝播を止める
	StopPropagation UpdateResult = 2 // コントロール内での以降の伝搬を停止する（親への戻り値としてはConsumed相当）
)

// ControlBaseを埋め込むとControlインターフェースを満たす
type Control interface {
	Update(t TouchInfo) UpdateResult
	Draw(screen *ebiten.Image)
	GetControlBase() *ControlBase
}

// ControlBaseはコントロールの基本機能を実装した構造体
// ・親となるコントロールを持つ
// ・座標とサイズを持つ(使わないこともあるがGUI用コントロールなので)
// ・Update/Drawのメソッドと、メソッド呼び出し時に実行する関数のスライスを持つ
type ControlBase struct {
	Control         Control
	Parent          Control
	X, Y            int
	Width, Height   int
	Visible         bool
	AutoResizable   bool
	updateFunctions []func(t TouchInfo) UpdateResult
	drawFunctions   []func(screen *ebiten.Image)
}

// ControlBase生成
func NewControlBase(c Control, x, y, w, h int) *ControlBase {
	return &ControlBase{Control: c, X: x, Y: y, Width: w, Height: h, Visible: true, AutoResizable: false}
}

// このコントロールのUpdateを呼び出した際に実行するUpdateFunctionを登録する
func (c *ControlBase) AddUpdateFunction(f func(t TouchInfo) UpdateResult) {
	c.updateFunctions = append(c.updateFunctions, f)
}

// 登録されたUpdateFunctionsを順次実行する
func (c *ControlBase) Update(t TouchInfo) UpdateResult {
	var r UpdateResult = NotConsumed

	if c.Visible {
		for i := len(c.updateFunctions) - 1; i >= 0; i-- {
			tmp := c.updateFunctions[i](t)

			// 消費された
			if tmp == Consumed {
				r = Consumed
			}

			// 伝搬即時停止
			if tmp == StopPropagation {
				r = Consumed
				t = nil
			}
		}
	}

	return r
}

// このコントロールのDrawを呼び出した際に実行するDrawFunctionを登録する
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

// 親Controlの座標を取得する
func (c *ControlBase) GetOwnerPos() (int, int) {
	p := c.Parent
	if p == nil {
		return 0, 0
	} else {
		pc := p.GetControlBase()
		x, y := pc.GetOwnerPos()
		return x + pc.X, y + pc.Y
	}
}

// ControlBaseのポインタを返す。Controlインターフェースのメソッドの一つ。
// Controlとしての基本機能はControlBaseに定義することでControlインターフェースをなるべくシンプルにしておきたい。
func (c *ControlBase) GetControlBase() *ControlBase {
	return c
}
