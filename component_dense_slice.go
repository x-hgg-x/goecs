package ecs

import "github.com/yourbasic/bit"

// DenseSliceComponent uses a 2-way redirection table between entities and components, allowing to leave no gaps in the data slice.
// Consumes less memory if component is a big struct.
type DenseSliceComponent struct {
	component
	data     []interface{}
	dataID   []int
	entityID []Entity
}

// Get returns data corresponding to entity
func (c *DenseSliceComponent) Get(entity Entity) interface{} {
	if 0 <= int(entity) && int(entity) < len(c.dataID) && c.dataID[entity] != -1 {
		return c.data[c.dataID[entity]]
	}
	return nil
}

// Set sets data corresponding to entity, or does nothing if the entity does not have the component
func (c *DenseSliceComponent) Set(entity Entity, data interface{}) {
	if entity.HasComponent(c) {
		deltaLen := int(entity) + 1 - len(c.dataID)
		for iSlice := 0; iSlice < deltaLen; iSlice++ {
			c.dataID = append(c.dataID, -1)
		}

		c.dataID[entity] = len(c.data)
		c.entityID = append(c.entityID, entity)
		c.data = append(c.data, data)
	}
}

func (c *DenseSliceComponent) _Remove(entity Entity) {
	if !(0 <= int(entity) && int(entity) < len(c.dataID) && c.dataID[entity] != -1) {
		return
	}

	// Delete component and replace it by the last element of the data slice
	id := c.dataID[entity]
	c.dataID[c.entityID[len(c.entityID)-1]] = id
	c.dataID[entity] = -1

	c.entityID[id] = c.entityID[len(c.entityID)-1]
	c.entityID = c.entityID[:len(c.entityID)-1]

	c.data[id] = c.data[len(c.data)-1]
	c.data = c.data[:len(c.data)-1]
}

func (c *DenseSliceComponent) _Reset() {
	c.data, c.dataID, c.entityID = nil, nil, nil
	c.component.tag = bit.Set{}
}
