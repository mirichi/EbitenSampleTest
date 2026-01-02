package systems

import (
	"MyProject/objects"
)

// Collider is an interface for entities that can be collided with.
type Collider interface {
	objects.CollisionTester
	IsDead() bool
}

// CollidableGroup is a pair of a Container and the type of items it holds (implied by usage).
type CollidableGroup struct {
	Container *objects.Container
}

// CollisionManager handles collision detection between groups.
type CollisionManager struct{}

// NewCollisionManager creates a new CollisionManager.
func NewCollisionManager() *CollisionManager {
	return &CollisionManager{}
}

// CheckCollision checks for collisions between an entity and a group of entities.
// onHit is called when a collision occurs, passing the entity from the group.
// It returns true if a collision occurred.
func (cm *CollisionManager) CheckCollision(entity Collider, group *objects.Container, onHit func(target interface{})) bool {
	if entity.IsDead() {
		return false
	}

	for _, child := range group.Children {
		target, ok := child.(Collider)
		if !ok || target.IsDead() {
			continue
		}

		if entity.Test(target) {
			onHit(target)
			return true
		}
	}
	return false
}

// CheckGroupCollision checks for collisions between two groups of entities.
// onHit is called for each collision pair (entityFromA, entityFromB).
func (cm *CollisionManager) CheckGroupCollision(groupA *objects.Container, groupB *objects.Container, onHit func(a, b interface{})) {
	for _, childA := range groupA.Children {
		entityA, okA := childA.(Collider)
		if !okA || entityA.IsDead() {
			continue
		}

		for _, childB := range groupB.Children {
			entityB, okB := childB.(Collider)
			if !okB || entityB.IsDead() {
				continue
			}

			if entityA.Test(entityB) {
				onHit(entityA, entityB)
			}
		}
	}
}
