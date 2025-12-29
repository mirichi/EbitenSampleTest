// entity.go
package parts

import (
	"MyProject/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	HandleInput(t input.Touch) bool
	Update()
	Draw(screen *ebiten.Image)
	GetGlobalPos() (float64, float64)
	GetEntityBase() *EntityBase
}

type EntityBase struct {
	Entity  Entity
	Parent  Entity
	X, Y    float64
	Visible bool

	// 実行タイミングをbeforeとafterで制御できる
	// 基本は真ん中で実行だが、前処理や後処理をしたいときに活用する

	// 入力処理
	beforeHandleInputFunctions []func(t input.Touch) bool // Groupingの配下の処理などが始まる前の前処理
	handleInputFunctions       []func(t input.Touch) bool // メインの入力処理(Groupingの配下含む)
	afterHandleInputFunctions  []func(t input.Touch) bool // Groupingの配下の処理などが終わってからの後処理

	// 更新処理
	beforeUpdateFunctions []func() // Groupingの配下の処理などが始まる前の前処理
	updateFunctions       []func() // メインの更新処理(Groupingの配下含む)
	afterUpdateFunctions  []func() // Groupingの配下の処理などが終わってからの後処理

	// 描画処理
	beforeDrawFunctions []func(screen *ebiten.Image) // 背景的なものの描画
	drawFunctions       []func(screen *ebiten.Image) // メインの描画(Groupingの配下含む)
	afterDrawFunctions  []func(screen *ebiten.Image) // 文字など前面の描画
}

func NewEntityBase(e Entity, x, y float64) *EntityBase {
	eb := &EntityBase{}
	eb.InitEntityBase(e, x, y)
	return eb
}

func (eb *EntityBase) InitEntityBase(e Entity, x, y float64) {
	eb.Entity = e
	eb.X = x
	eb.Y = y
	eb.Visible = true
}

// AddHandleInputFunctionは、HandleInputメソッド呼び出し時に実行される関数を登録する
// 登録された関数は後に追加されたものから順に実行される
func (eb *EntityBase) AddBeforeHandleInputFunction(f func(t input.Touch) bool) {
	eb.beforeHandleInputFunctions = append(eb.beforeHandleInputFunctions, f)
}
func (eb *EntityBase) AddHandleInputFunction(f func(t input.Touch) bool) {
	eb.handleInputFunctions = append(eb.handleInputFunctions, f)
}
func (eb *EntityBase) AddAfterHandleInputFunction(f func(t input.Touch) bool) {
	eb.afterHandleInputFunctions = append(eb.afterHandleInputFunctions, f)
}

// 登録されたHandleInputFunctionsを後ろから順次実行する
func (eb *EntityBase) HandleInput(t input.Touch) bool {
	handle := false

	if !eb.Visible {
		t = nil
	}

	for i := len(eb.beforeHandleInputFunctions) - 1; i >= 0; i-- {
		if eb.beforeHandleInputFunctions[i](t) {
			t = nil
			handle = true
		}
	}
	for i := len(eb.handleInputFunctions) - 1; i >= 0; i-- {
		if eb.handleInputFunctions[i](t) {
			t = nil
			handle = true
		}
	}
	for i := len(eb.afterHandleInputFunctions) - 1; i >= 0; i-- {
		if eb.afterHandleInputFunctions[i](t) {
			t = nil
			handle = true
		}
	}

	return handle
}

// AddUpdateFunctionは、Updateメソッド呼び出し時に実行される関数を登録する
// 登録された関数は後に追加されたものから順に実行される
func (eb *EntityBase) AddBeforeUpdateFunction(f func()) {
	eb.beforeUpdateFunctions = append(eb.beforeUpdateFunctions, f)
}
func (eb *EntityBase) AddUpdateFunction(f func()) {
	eb.updateFunctions = append(eb.updateFunctions, f)
}
func (eb *EntityBase) AddAfterUpdateFunction(f func()) {
	eb.afterUpdateFunctions = append(eb.afterUpdateFunctions, f)
}

// 登録されたUpdateFunctionsを後ろから順次実行する
func (eb *EntityBase) Update() {
	for i := len(eb.beforeUpdateFunctions) - 1; i >= 0; i-- {
		eb.beforeUpdateFunctions[i]()
	}
	for i := len(eb.updateFunctions) - 1; i >= 0; i-- {
		eb.updateFunctions[i]()
	}
	for i := len(eb.afterUpdateFunctions) - 1; i >= 0; i-- {
		eb.afterUpdateFunctions[i]()
	}
}

// AddDrawFunctionは、Drawメソッド呼び出し時に実行される関数を登録する
// 登録された関数は登録順に実行される
func (eb *EntityBase) AddBeforeDrawFunction(f func(screen *ebiten.Image)) {
	eb.beforeDrawFunctions = append(eb.beforeDrawFunctions, f)
}
func (eb *EntityBase) AddDrawFunction(f func(screen *ebiten.Image)) {
	eb.drawFunctions = append(eb.drawFunctions, f)
}
func (eb *EntityBase) AddAfterDrawFunction(f func(screen *ebiten.Image)) {
	eb.afterDrawFunctions = append(eb.afterDrawFunctions, f)
}

// 登録されたDrawFunctionsを順次実行する
func (eb *EntityBase) Draw(screen *ebiten.Image) {
	if eb.Visible {
		for _, f := range eb.beforeDrawFunctions {
			f(screen)
		}
		for _, f := range eb.drawFunctions {
			f(screen)
		}
		for _, f := range eb.afterDrawFunctions {
			f(screen)
		}
	}
}

// getGlobalPosは、画面全体（ルートコントロール）からの絶対座標を取得
// 親コントロールの座標を再帰的に加算して算出する
func (eb *EntityBase) GetGlobalPos() (float64, float64) {
	p := eb.Parent
	if p == nil {
		return eb.X, eb.Y
	} else {
		x, y := p.GetGlobalPos()
		return x + eb.X, y + eb.Y
	}
}

// EntityBaseのポインタを返す。Entityインターフェースのメソッドの一つ
// Entityとしての基本機能はEntityBaseに定義することでEntityインターフェースをなるべくシンプルにしておきたい
func (eb *EntityBase) GetEntityBase() *EntityBase {
	return eb
}
