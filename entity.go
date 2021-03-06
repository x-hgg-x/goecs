package ecs

// Entity is an index
type Entity int

// AddComponent adds entity for component
func (entity Entity) AddComponent(component DataComponent, data interface{}) Entity {
	component._Tag().Add(int(entity))
	component.Set(entity, data)
	component._Manager().entities.Set(component._Manager().getEntities())
	return entity
}

// RemoveComponent removes entity for component
func (entity Entity) RemoveComponent(component DataComponent) Entity {
	component._Tag().Delete(int(entity))
	component._Remove(entity)
	component._Manager().entities.Set(component._Manager().getEntities())
	return entity
}

// HasComponent checks if component has entity
func (entity Entity) HasComponent(component DataComponent) bool {
	return component._Tag().Contains(int(entity))
}
