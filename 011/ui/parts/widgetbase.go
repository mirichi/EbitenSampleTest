package parts

import (
	"MyProject/ui/input"

	"github.com/hajimehoshi/ebiten/v2"
)

// WidgetBaseを埋め込むとWidgetインターフェースを満たす
type Widget interface {
	HandleInput(t input.Touch) bool
	Update()
	Draw(screen *ebiten.Image)
	GetWidgetBase() *WidgetBase
}

// WidgetBaseはWidgetインターフェース基本機能を実装した構造体
// 親ウィジェットの参照、座標、サイズ、可視性などの共通プロパティを管理する
// また、HandleInput/Update/Drawの各フェーズで実行される関数リストを保持している
type WidgetBase struct {
	Widget               Widget
	Parent               Widget
	X, Y                 int
	Width, Height        int
	Visible              bool
	AutoResizable        bool
	handleInputFunctions []func(t input.Touch) bool
	updateFunctions      []func()
	drawFunctions        []func(screen *ebiten.Image)
}

func NewWidgetBase(c Widget, x, y, w, h int) *WidgetBase {
	cb := &WidgetBase{}
	cb.InitWidgetBase(c, x, y, w, h)
	return cb
}

func (cb *WidgetBase) InitWidgetBase(c Widget, x, y, w, h int) {
	cb.Widget = c
	cb.X = x
	cb.Y = y
	cb.Width = w
	cb.Height = h
	cb.Visible = true
	cb.AutoResizable = false
}

// AddHandleInputFunctionは、HandleInputメソッド呼び出し時に実行される関数を登録する
// 登録された関数は後に追加されたものから順に実行される
func (cb *WidgetBase) AddHandleInputFunction(f func(t input.Touch) bool) {
	cb.handleInputFunctions = append(cb.handleInputFunctions, f)
}

// 登録されたHandleInputFunctionsを順次実行する
func (cb *WidgetBase) HandleInput(t input.Touch) bool {
	if cb.Visible && t != nil {
		for i := len(cb.handleInputFunctions) - 1; i >= 0; i-- {
			if cb.handleInputFunctions[i](t) {
				return true
			}
		}
	}

	return false
}

// AddUpdateFunctionは、Updateメソッド呼び出し時に実行される関数を登録する
// 登録された関数は後に追加されたものから順に実行される
func (cb *WidgetBase) AddUpdateFunction(f func()) {
	cb.updateFunctions = append(cb.updateFunctions, f)
}

// 登録されたUpdateFunctionsを順次実行する
func (cb *WidgetBase) Update() {
	for i := len(cb.updateFunctions) - 1; i >= 0; i-- {
		cb.updateFunctions[i]()
	}
}

// AddDrawFunctionは、Drawメソッド呼び出し時に実行される関数を登録する
// 登録された関数は登録順に実行される
func (cb *WidgetBase) AddDrawFunction(f func(screen *ebiten.Image)) {
	cb.drawFunctions = append(cb.drawFunctions, f)
}

// 登録されたDrawFunctionsを順次実行する
func (cb *WidgetBase) Draw(screen *ebiten.Image) {
	if cb.Visible {
		for _, f := range cb.drawFunctions {
			f(screen)
		}
	}
}

// getGlobalPosは、画面全体（ルートコントロール）からの絶対座標を取得
// 親コントロールの座標を再帰的に加算して算出する
func (cb *WidgetBase) GetGlobalPos() (int, int) {
	p := cb.Parent
	if p == nil {
		return cb.X, cb.Y
	} else {
		pcb := p.GetWidgetBase()
		x, y := pcb.GetGlobalPos()
		return x + cb.X, y + cb.Y
	}
}

// WidgetBaseのポインタを返す。Widgetインターフェースのメソッドの一つ
// Widgetとしての基本機能はWidgetBaseに定義することでWidgetインターフェースをなるべくシンプルにしておきたい
func (cb *WidgetBase) GetWidgetBase() *WidgetBase {
	return cb
}
