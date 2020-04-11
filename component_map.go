package ecs

import "github.com/yourbasic/bit"

// MapComponent is a data storage
type MapComponent struct {
	component
	data map[Entity]interface{}
}

// Get returns data corresponding to entity
func (c *MapComponent) Get(entity Entity) interface{} {
	if data, ok := c.data[entity]; ok {
		return data
	}
	return nil
}

// Set sets data corresponding to entity, or does nothing if the entity does not have the component
func (c *MapComponent) Set(entity Entity, data interface{}) {
	if entity.HasComponent(c) {
		c.data[entity] = data
	}
}

func (c *MapComponent) _Remove(entity Entity) {
	delete(c.data, entity)
}

func (c *MapComponent) _Reset() {
	c.data = make(map[Entity]interface{})
	c.component.tag = bit.Set{}
}
