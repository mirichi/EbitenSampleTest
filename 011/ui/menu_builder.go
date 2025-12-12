package ui

// MenuBuilder はPopupMenuを構築するためのビルダー
type MenuBuilder struct {
	x, y  int
	items []MenuItem
}

// NewMenuBuilder は新しいMenuBuilderを作成する
func NewMenuBuilder(x, y int) *MenuBuilder {
	return &MenuBuilder{
		x:     x,
		y:     y,
		items: make([]MenuItem, 0),
	}
}

// AddItem はアクション付きの項目を追加する
func (b *MenuBuilder) AddItem(text string, action func()) *MenuBuilder {
	b.items = append(b.items, MenuItem{
		Text:   text,
		Action: action,
	})
	return b
}

// AddDisabledItem は無効化された項目を追加する
func (b *MenuBuilder) AddDisabledItem(text string) *MenuBuilder {
	b.items = append(b.items, MenuItem{
		Text:     text,
		Disabled: true,
	})
	return b
}

// AddSeparator はセパレータを追加する
func (b *MenuBuilder) AddSeparator() *MenuBuilder {
	b.items = append(b.items, MenuItem{
		IsSeparator: true,
	})
	return b
}

func (b *MenuBuilder) AddSubMemu(text string, m *PopupMenu) *MenuBuilder {
	b.items = append(b.items, MenuItem{
		Text:    text,
		SubMenu: m,
	})
	return b
}

// Build はPopupMenuを生成する
func (b *MenuBuilder) Build() *PopupMenu {
	return NewPopupMenu(b.x, b.y, b.items)
}
