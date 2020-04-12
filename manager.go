package ecs

import "github.com/yourbasic/bit"

// Manager manages components and entities
type Manager struct {
	currentEntityIndex int
	components         []DataComponent
	entities           *bit.Set
}

// NewManager creates a new manager
func NewManager() *Manager {
	return &Manager{entities: bit.New()}
}

// NewNullComponent creates a new NullComponent
func (manager *Manager) NewNullComponent() *NullComponent {
	component := &NullComponent{component: component{manager: manager}}
	manager.components = append(manager.components, component)
	return component
}

// NewSliceComponent creates a new SliceComponent
func (manager *Manager) NewSliceComponent() *SliceComponent {
	component := &SliceComponent{component: component{manager: manager}}
	manager.components = append(manager.components, component)
	return component
}

// NewDenseSliceComponent creates a new DenseSliceComponent
func (manager *Manager) NewDenseSliceComponent() *DenseSliceComponent {
	component := &DenseSliceComponent{component: component{manager: manager}}
	manager.components = append(manager.components, component)
	return component
}

// NewMapComponent creates a new MapComponent
func (manager *Manager) NewMapComponent() *MapComponent {
	component := &MapComponent{data: make(map[Entity]interface{}), component: component{manager: manager}}
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
		manager.components[iComponent]._Reset()
	}
	// Reset current entity index and entity list
	manager.currentEntityIndex = 0
	manager.entities = bit.New()
}

// Join returns tag describing intersection of components
func (manager *Manager) Join(components ...joinable) *bit.Set {
	tag := bit.New().Set(manager.entities)
	for _, component := range components {
		tag = component._Join(tag)
	}
	return tag
}

// Maintain reorders component storage to eliminate gaps and reduce memory usage.
// It invalidates all previous entities.
func (manager *Manager) Maintain(minEntities int, maxFillRatio float64) {
	if manager.entities.Empty() {
		return
	}
	if manager.currentEntityIndex < minEntities || float64(manager.entities.Size())/float64(manager.currentEntityIndex) > maxFillRatio {
		return
	}

	// Get entities without components
	diffTag := bit.New().AddRange(0, manager.entities.Max()+1)
	diffTag.SetXor(diffTag, manager.entities)

	diffTag.Visit(func(n int) (skip bool) {
		if n >= manager.entities.Max() {
			return true
		}

		// Copy the components of the last entity to an empty entity, then delete the last entity
		maxEntity := Entity(manager.entities.Max())
		for _, component := range manager.components {
			if maxEntity.HasComponent(component) {
				Entity(n).AddComponent(component, component.Get(maxEntity))
				maxEntity.RemoveComponent(component)
			}
		}
		return false
	})

	manager.currentEntityIndex = manager.entities.Max() + 1
}

// Get all entities with at least one component
func (manager *Manager) getEntities() *bit.Set {
	tag := bit.New()
	for _, component := range manager.components {
		tag.SetOr(tag, component._Tag())
	}
	return tag
}
