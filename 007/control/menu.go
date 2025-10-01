package control

import (
	"image"
	"image/color"

	"myproject/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var emptyImage = ebiten.NewImage(3, 3)
var whitePixel = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

func init() {
	whitePixel.Fill(color.White)
}

// Menu
type Menu struct {
	ui.ControlBase // これを埋め込むとコントロールとして扱える
	Easing         func(x float64) float64
	From, To       int
	ElapsedFrames  int
	Duration       int
	EasingCompFunc func()
}

func NewMenu(x, y, w, h int, o ui.Control) *Menu {
	return &Menu{
		ControlBase: ui.ControlBase{X: x, Y: y, W: w, H: h, Owner: o},
	}
}

// 描画
func (m *Menu) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(m.X), float32(m.Y), float32(m.W), float32(m.H), color.RGBA{0x80, 0x80, 0x80, 0xf0}, false)

	var path vector.Path

	path.MoveTo(float32(m.X), float32(m.Y))
	path.LineTo(float32(m.X+m.W), float32(m.Y))
	path.LineTo(float32(m.X+m.W), float32(m.Y+m.H))
	path.LineTo(float32(m.X), float32(m.Y+m.H))
	path.Close()

	op := &vector.StrokeOptions{}
	op.Width = 10
	op.LineJoin = vector.LineJoinRound
	var vertices []ebiten.Vertex = []ebiten.Vertex{}
	var indices []uint16 = []uint16{}
	vertices, indices = path.AppendVerticesAndIndicesForStroke(vertices[:0], indices[:0], op)
	for i := range vertices {
		vertices[i].SrcX = 1
		vertices[i].SrcY = 1
		vertices[i].ColorR = 0xd0 / float32(0xff)
		vertices[i].ColorG = 0xd0 / float32(0xff)
		vertices[i].ColorB = 0xd0 / float32(0xff)
		vertices[i].ColorA = 1
	}
	op2 := &ebiten.DrawTrianglesOptions{}
	op2.AntiAlias = true
	op2.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(vertices, indices, whitePixel, op2)

	// 配下のコントロール描画
	for _, c := range m.Controls {
		c.Draw(screen)
	}
}

func (m *Menu) Update() {
	// イージング処理
	if m.Easing != nil {
		m.ElapsedFrames += 1
		t := m.Easing(float64(m.ElapsedFrames) / float64(m.Duration))
		m.X = m.From + int(float64(m.To-m.From)*t+0.5)

		if m.ElapsedFrames == m.Duration {
			m.Easing = nil
			if m.EasingCompFunc != nil {
				m.EasingCompFunc()
			}
		}
	}

	m.ControlBase.Update()
}

// メニュー画面(背景)
type MenuScreen struct {
	ui.ControlBase
	Running bool
}

// メニュー外のタップで消えるように全画面サイズのコントロールを作成する
func NewMenuScreen(x, y, w, h int, o ui.Control) *MenuScreen {
	ms := &MenuScreen{
		ControlBase: ui.ControlBase{X: x, Y: y, W: w, H: h, Owner: o},
		Running:     false,
	}

	// メニュー外タップでメニュー画面を消す処理
	ms.TapFunc = func() {
		m := ms.Controls[0].(*Menu)

		// アニメーション中のタップは無視
		if m.ElapsedFrames != m.Duration {
			return
		}
		m.Easing = func(x float64) float64 {
			return x * x * x * x * x
		}
		m.From = 100
		m.To = 640
		m.ElapsedFrames = 0
		m.Duration = 30
		m.EasingCompFunc = func() {
			ms.Running = false
		}
	}

	return ms
}

func (ms *MenuScreen) ProcessTouch(t ui.TouchInfo) bool {
	// メニュー画面起動中以外は処理しない
	if ms.Running {
		return ms.ControlBase.ProcessTouch(t)
	}

	return false
}

func (ms *MenuScreen) Draw(screen *ebiten.Image) {
	// メニュー画面起動中以外は処理しない
	if !ms.Running {
		return
	}

	vector.DrawFilledRect(screen, float32(ms.X), float32(ms.Y), float32(ms.W), float32(ms.H), color.RGBA{0x00, 0x00, 0x00, 0xb0}, false)

	// 配下のコントロール描画
	for _, c := range ms.Controls {
		c.Draw(screen)
	}
}

func (ms *MenuScreen) Update() {
	// メニュー画面起動中以外は処理しない
	if !ms.Running {
		return
	}

	ms.ControlBase.Update()
}

func (ms *MenuScreen) Start() {
	ms.Running = true

	// MenuScreen配下はMenu1個と決まっているので直接アクセスする
	m := ms.Controls[0].(*Menu)
	m.Easing = func(x float64) float64 {
		f := x - 1
		return f*f*f*f*f + 1
	}
	m.From = -440
	m.To = 100
	m.ElapsedFrames = 0
	m.Duration = 40
	m.EasingCompFunc = nil
}
