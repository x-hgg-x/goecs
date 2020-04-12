package ecs

import "github.com/yourbasic/bit"

// SliceComponent uses a slice for storing data.
// Useful for pointer components, for small struct components or for components mostly present in entities.
type SliceComponent struct {
	component
	data []interface{}
}

// Get returns data corresponding to entity
func (c *SliceComponent) Get(entity Entity) interface{} {
	if 0 <= int(entity) && int(entity) < len(c.data) {
		return c.data[entity]
	}
	return nil
}

// Set sets data corresponding to entity, or does nothing if the entity does not have the component
func (c *SliceComponent) Set(entity Entity, data interface{}) {
	// Check existing data
	if 0 <= int(entity) && int(entity) < len(c.data) && c.data[entity] != nil {
		c.data[entity] = data
		return
	}

	// Insert data
	if entity.HasComponent(c) {
		deltaLen := int(entity) + 1 - len(c.data)
		for iSlice := 0; iSlice < deltaLen; iSlice++ {
			c.data = append(c.data, nil)
		}
		c.data[entity] = data
	}
}

func (c *SliceComponent) _Remove(entity Entity) {
	if 0 <= int(entity) && int(entity) < len(c.data) {
		c.data[entity] = nil
	}
}

func (c *SliceComponent) _Reset() {
	c.data = nil
	c.component.tag = bit.Set{}
}
