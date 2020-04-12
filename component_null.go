package ecs

import "github.com/yourbasic/bit"

// NullComponent contains no data and works as a simple flag.
type NullComponent struct {
	component
}

// Get returns data corresponding to entity
func (c *NullComponent) Get(entity Entity) interface{} {
	return nil
}

// Set sets data corresponding to entity, or does nothing if the entity does not have the component
func (c *NullComponent) Set(entity Entity, data interface{}) {}

func (c *NullComponent) _Remove(entity Entity) {}

func (c *NullComponent) _Reset() {
	c.component.tag = bit.Set{}
}
