package ui

type Focusable struct {
	Control Control
	Focused bool

	OnFocus func()
	OnBlur  func()
}

type FocusableInterface interface {
	Focus()
	Blur()
}

var FocusedControl FocusableInterface
var FlameFocus bool

func NewFocusable(c Control) *Focusable {
	f := &Focusable{Control: c}
	return f
}

func (f *Focusable) Focus() {
	if FocusedControl != f.Control.(FocusableInterface) {
		if FocusedControl != nil {
			FocusedControl.Blur()
		}
		f.Focused = true
		FocusedControl = f.Control.(FocusableInterface)
		if f.OnFocus != nil {
			f.OnFocus()
		}
	}
	FlameFocus = true
}

func (f *Focusable) Blur() {
	f.Focused = false
	FocusedControl = nil
	if f.OnBlur != nil {
		f.OnBlur()
	}
}
