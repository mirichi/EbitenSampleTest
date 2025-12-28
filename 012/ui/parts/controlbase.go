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
	GetGlobalPos() (int, int)
	GetControlBase() *ControlBase
}

// ControlBaseはControlインターフェース基本機能を実装した構造体
// 親Controlの参照、座標、サイズ、可視性などの共通プロパティを管理する
// また、HandleInput/Update/Drawの各フェーズで実行される関数リストを保持している
type ControlBase struct {
	EntityBase

	Control         Control
	Width, Height   int
	FlexGrow        float64
	FlexShrink      float64
	FlexBasisWidth  int
	FlexBasisHeight int
}

func NewControlBase(c Control, x, y, w, h int) *ControlBase {
	cb := &ControlBase{}
	cb.InitControlBase(c, x, y, w, h)
	return cb
}

func (cb *ControlBase) InitControlBase(c Control, x, y, w, h int) {
	cb.EntityBase.InitEntityBase(c.(Entity), x, y)
	cb.Control = c
	cb.Width = w
	cb.Height = h
	cb.FlexGrow = 0
	cb.FlexShrink = 0
	cb.FlexBasisWidth = w
	cb.FlexBasisHeight = h
}

func (cb *ControlBase) GetSize() (int, int) {
	return cb.Width, cb.Height
}

// ControlBaseのポインタを返す。Controlインターフェースのメソッドの一つ
// Controlとしての基本機能はControlBaseに定義することでControlインターフェースをなるべくシンプルにしておきたい
func (cb *ControlBase) GetControlBase() *ControlBase {
	return cb
}
