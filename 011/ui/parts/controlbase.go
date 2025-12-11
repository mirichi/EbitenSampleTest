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

// ControlBaseはControlインターフェース基本機能を実装した構造体
// 親Controlの参照、座標、サイズ、可視性などの共通プロパティを管理する
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
	OnBeforeHandleInput  func()
	OnAfterHandleInput   func()
	OnBeforeUpdate       func()
	OnAfterUpdate        func()
}

func NewControlBase(c Control, x, y, w, h int) *ControlBase {
	cb := &ControlBase{}
	cb.InitControlBase(c, x, y, w, h)
	return cb
}

func (cb *ControlBase) InitControlBase(c Control, x, y, w, h int) {
	cb.Control = c
	cb.X = x
	cb.Y = y
	cb.Width = w
	cb.Height = h
	cb.Visible = true
	cb.AutoResizable = false
}

// AddHandleInputFunctionは、HandleInputメソッド呼び出し時に実行される関数を登録する
// 登録された関数は後に追加されたものから順に実行される
func (cb *ControlBase) AddHandleInputFunction(f func(t input.Touch) bool) {
	cb.handleInputFunctions = append(cb.handleInputFunctions, f)
}

// 登録されたHandleInputFunctionsを順次実行する
func (cb *ControlBase) HandleInput(t input.Touch) bool {
	handle := false
	if cb.OnBeforeHandleInput != nil {
		cb.OnBeforeHandleInput()
	}

	if !cb.Visible {
		t = nil
	}

	for i := len(cb.handleInputFunctions) - 1; i >= 0; i-- {
		if cb.handleInputFunctions[i](t) {
			t = nil
			handle = true
		}
	}

	if cb.OnAfterHandleInput != nil {
		cb.OnAfterHandleInput()
	}

	return handle
}

// AddUpdateFunctionは、Updateメソッド呼び出し時に実行される関数を登録する
// 登録された関数は後に追加されたものから順に実行される
func (cb *ControlBase) AddUpdateFunction(f func()) {
	cb.updateFunctions = append(cb.updateFunctions, f)
}

// 登録されたUpdateFunctionsを順次実行する
func (cb *ControlBase) Update() {
	if cb.OnBeforeUpdate != nil {
		cb.OnBeforeUpdate()
	}

	for i := len(cb.updateFunctions) - 1; i >= 0; i-- {
		cb.updateFunctions[i]()
	}

	if cb.OnAfterUpdate != nil {
		cb.OnAfterUpdate()
	}
}

// AddDrawFunctionは、Drawメソッド呼び出し時に実行される関数を登録する
// 登録された関数は登録順に実行される
func (cb *ControlBase) AddDrawFunction(f func(screen *ebiten.Image)) {
	cb.drawFunctions = append(cb.drawFunctions, f)
}

// 登録されたDrawFunctionsを順次実行する
func (cb *ControlBase) Draw(screen *ebiten.Image) {
	if cb.Visible {
		for _, f := range cb.drawFunctions {
			f(screen)
		}
	}
}

// getGlobalPosは、画面全体（ルートコントロール）からの絶対座標を取得
// 親コントロールの座標を再帰的に加算して算出する
func (cb *ControlBase) GetGlobalPos() (int, int) {
	p := cb.Parent
	if p == nil {
		return cb.X, cb.Y
	} else {
		pcb := p.GetControlBase()
		x, y := pcb.GetGlobalPos()
		return x + cb.X, y + cb.Y
	}
}

// ControlBaseのポインタを返す。Controlインターフェースのメソッドの一つ
// Controlとしての基本機能はControlBaseに定義することでControlインターフェースをなるべくシンプルにしておきたい
func (cb *ControlBase) GetControlBase() *ControlBase {
	return cb
}
