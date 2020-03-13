package ecs

import "github.com/yourbasic/bit"

// Manager manages components and entities
type Manager struct {
	currentEntityIndex int
	components         []*Component
	entities           *bit.Set
}

// NewManager creates a new manager
func NewManager() *Manager {
	return &Manager{entities: bit.New()}
}

// NewComponent creates a new component
func (manager *Manager) NewComponent() *Component {
	component := &Component{data: make(map[Entity]interface{}), manager: manager}
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
		if entity.HasComponent(component) {
			entity.RemoveComponent(component)
		}
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
	for iComponent := range manager.components {
		*manager.components[iComponent] = Component{data: make(map[Entity]interface{}), manager: manager}
	}
	// Reset current entity index and entity list
	manager.currentEntityIndex = 0
	manager.entities = bit.New()
}

// Join returns tag describing intersection of components
func (manager *Manager) Join(components ...component) *bit.Set {
	tag := bit.New().Set(manager.entities)
	for _, component := range components {
		tag = component._Join(tag)
	}
	return tag
}

// Get all entities with at least one component
func (manager *Manager) getEntities() *bit.Set {
	tag := bit.New()
	for _, component := range manager.components {
		tag.SetOr(tag, &component.tag)
	}
	return tag
}
