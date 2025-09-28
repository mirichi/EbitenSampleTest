package ui

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TextAdjust int

const (
	AdjustLeft   TextAdjust = 0
	AdjustCenter TextAdjust = 1
	AdjustRight  TextAdjust = 2
)

var MplusFaceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	MplusFaceSource = s
}

type Control interface {
	ProcessTouch(t TouchInfo) bool
	Update()
	Draw(screen *ebiten.Image)
	GetGlobalPos() (int, int)
}

type ControlBase struct {
	X, Y      int
	W, H      int
	TouchInfo TouchInfo
	Owner     Control
	Controls  []Control
	TapFunc   func()
	SlideFunc func()
}

func NewControlBase(x, y, w, h int, o Control) *ControlBase {
	return &ControlBase{
		X:     x,
		Y:     y,
		W:     w,
		H:     h,
		Owner: o,
	}
}

func (cb *ControlBase) CheckPoint(t TouchInfo) bool {
	x, y := t.Pos()
	ox, oy := cb.GetGlobalPos()
	return x >= ox && x < ox+cb.W && y >= oy && y < oy+cb.H
}

// タッチした情報を保持する
func (cb *ControlBase) ProcessTouch(t TouchInfo) bool {
	// 既にタッチ中
	if cb.TouchInfo != nil {
		return true
	}

	// 自分の範囲がタッチされた
	if cb.CheckPoint(t) {
		for _, c := range cb.Controls {
			if c.ProcessTouch(t) {
				// 配下の何かがタッチされた
				return true
			}
		}

		// 配下の何かがタッチされていなければ自分がタッチされた
		cb.TouchInfo = t
		return true
	}

	// 自分の範囲外だった
	return false
}

// タップとスライドの判定とメソッド呼び出し
func (cb *ControlBase) Update() {
	// 自分がタッチ中
	if cb.TouchInfo != nil {
		// タッチが離された
		if !cb.TouchInfo.IsPressed() {
			// タップイベント
			if cb.TapFunc != nil {
				cb.TapFunc()
			}
			cb.TouchInfo = nil // タッチ中を解除
		} else if !cb.CheckPoint(cb.TouchInfo) {
			// ヒットしていない場合はタッチ中を解除
			cb.TouchInfo = nil
		} else {
			// タッチが移動していたらスライドイベント
			oldx, oldy := cb.TouchInfo.OldPos()
			x, y := cb.TouchInfo.Pos()
			if x != oldx || y != oldy {
				if cb.SlideFunc != nil {
					cb.SlideFunc()
				}
			}
		}
	}

	// 配下の処理
	for _, c := range cb.Controls {
		c.Update()
	}
}

// 最低限の枠だけ描画する機能を実装しておく
func (cb *ControlBase) Draw(screen *ebiten.Image) {
	ox, oy := cb.GetGlobalPos()
	vector.StrokeRect(screen, float32(ox), float32(oy), float32(cb.W), float32(cb.H), 3, color.RGBA{0xff, 0xff, 0xff, 0xff}, false)

	// 配下の処理
	for _, c := range cb.Controls {
		c.Draw(screen)
	}
}

// 描画や判定に用いるスクリーン座標を算出する
func (cb *ControlBase) GetGlobalPos() (int, int) {
	if cb.Owner == nil {
		return cb.X, cb.Y
	}

	ox, oy := cb.Owner.GetGlobalPos()
	return ox + cb.X, oy + cb.Y
}
