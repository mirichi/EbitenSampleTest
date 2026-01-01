package objects

import "github.com/hajimehoshi/ebiten/v2"

// Scene はゲームの各シーンが実装すべきインターフェース
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// SceneManager はシーンの遷移を管理する
type SceneManager struct {
	currentScene Scene
	nextScene    Scene
}

// NewSceneManager は新しいSceneManagerを作成する
func NewSceneManager(initialScene Scene) *SceneManager {
	return &SceneManager{
		currentScene: initialScene,
	}
}

// Update は現在のシーンを更新し、必要ならシーン遷移を行う
func (sm *SceneManager) Update() error {
	if sm.nextScene != nil {
		sm.currentScene = sm.nextScene
		sm.nextScene = nil
	}

	if sm.currentScene != nil {
		return sm.currentScene.Update()
	}
	return nil
}

// Draw は現在のシーンを描画する
func (sm *SceneManager) Draw(screen *ebiten.Image) {
	if sm.currentScene != nil {
		sm.currentScene.Draw(screen)
	}
}

// ChangeScene は次の更新タイミングでシーンを切り替える
func (sm *SceneManager) ChangeScene(next Scene) {
	sm.nextScene = next
}

// GetCurrentScene は現在のシーンを返す（型アサーション用）
func (sm *SceneManager) GetCurrentScene() Scene {
	return sm.currentScene
}
