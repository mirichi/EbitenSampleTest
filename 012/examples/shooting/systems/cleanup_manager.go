package systems

import (
	"MyProject/objects"
)

// Cleanable is an interface for entities that can be cleaned up (removed).
type Cleanable interface {
	IsDead() bool
}

// CleanupManager handles the removal of dead entities from containers.
type CleanupManager struct{}

// NewCleanupManager creates a new CleanupManager.
func NewCleanupManager() *CleanupManager {
	return &CleanupManager{}
}

// Cleanup removes dead entities from the specified container.
func (cm *CleanupManager) Cleanup(container *objects.Container) {
	children := container.Children
	// Iterate backwards to safely remove elements
	for i := len(children) - 1; i >= 0; i-- {
		child := children[i]
		if cleanable, ok := child.(Cleanable); ok {
			if cleanable.IsDead() {
				container.RemoveChild(child)
			}
		}
	}
}
