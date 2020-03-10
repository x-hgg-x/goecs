package ecs

import "github.com/yourbasic/bit"

// Manager manages components and entities
type Manager struct {
	currentEntityIndex int
	components         []*Component
}

// NewComponent creates a new component
func (manager *Manager) NewComponent() *Component {
	component := &Component{
		data: make(map[Entity]interface{}),
	}

	manager.components = append(manager.components, component)
	return component
}

// NewEntity creates a new entity
func (manager *Manager) NewEntity() Entity {
	manager.currentEntityIndex++
	return Entity(manager.currentEntityIndex - 1)
}

// DeleteEntity removes entity for all associated components
func (manager *Manager) DeleteEntity(entity Entity) {
	for _, component := range manager.components {
		entity.RemoveComponent(component)
	}
}

// DeleteEntities removes entities for all associated components
func (manager *Manager) DeleteEntities(entities ...Entity) {
	for _, entity := range entities {
		manager.DeleteEntity(entity)
	}
}

// DeleteAllEntities removes all entities for all components and reset current entity index
func (manager *Manager) DeleteAllEntities() {
	for _, component := range manager.components {
		component.tag = bit.Set{}
	}
	// Reset current entity index
	manager.currentEntityIndex = 0
}
