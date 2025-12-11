package ui

import (
	"MyProject/ui/input"
	"MyProject/ui/parts"
)

// PopupManager はポップアップウィンドウの管理を行う
// 将来的にはMainScreenやWindowごとにインスタンスを持つことを想定
type PopupManager struct {
	container *popupContainer
}

// NewPopupManager は新しいPopupManagerを作成する
func NewPopupManager() *PopupManager {
	return &PopupManager{
		// 既存のPopupContainerの仕組みを利用するが、
		// 将来的にはここをローカルなコンテナ生成に置き換える
		container: NewPopupContainer(),
	}
}

// Update は管理下のポップアップの状態を更新する
func (pm *PopupManager) Update() {
	if pm.container != nil {
		pm.container.Update()
	}
}

// HandleInput は入力イベントを処理する
func (pm *PopupManager) HandleInput(t input.Touch) bool {
	if pm.container != nil {
		return pm.container.HandleInput(t)
	}
	return false
}

// ShowMenu は指定されたメニューを表示する
func (pm *PopupManager) ShowMenu(menu *PopupMenu) {
	if pm.container == nil {
		return
	}
	// 既存のメニューを閉じる（単純化のため現在は1つのみ表示）
	pm.CloseAll()

	// コンテナに追加
	pm.container.AddChild(menu)

	// Menu自身にもManagerへの参照を持たせると閉じる処理などがスムーズになるが
	// 現状はCallbackで対応
	menu.OnClose = func() {
		pm.CloseAll()
	}
}

// CloseAll は全てのポップアップを閉じる
func (pm *PopupManager) CloseAll() {
	if pm.container != nil {
		pm.container.Close()
	}
}

// Container は描画対象となるContainerを返す
// MainScreenなどはこれをAddChildして描画ツリーに組み込む
func (pm *PopupManager) Container() parts.Control {
	return pm.container
}
