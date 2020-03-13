package ecs

// Entity is an index
type Entity int

// AddComponent adds entity for component
func (entity Entity) AddComponent(component *Component, data interface{}) Entity {
	component.tag.Add(int(entity))
	component.data[entity] = data
	component.manager.entities.Set(component.manager.getEntities())
	return entity
}

// RemoveComponent removes entity for component
func (entity Entity) RemoveComponent(component *Component) Entity {
	component.tag.Delete(int(entity))
	delete(component.data, entity)
	component.manager.entities.Set(component.manager.getEntities())
	return entity
}

// HasComponent checks if component has entity
func (entity Entity) HasComponent(component *Component) bool {
	return component.tag.Contains(int(entity))
}
