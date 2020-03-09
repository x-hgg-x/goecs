package ecs

import "github.com/yourbasic/bit"

type component interface {
	_Tag() *bit.Set
	_Join(*bit.Set) *bit.Set
}

// Join returns tag describing intersection of components
func Join(components ...component) *bit.Set {
	tag := &bit.Set{}
	if len(components) > 0 {
		tag.Set(components[0]._Tag())
		for _, component := range components[1:] {
			tag = component._Join(tag)
		}
	}
	return tag
}

// Component is a data storage
type Component struct {
	tag  bit.Set
	data map[Entity]interface{}
}

// Get returns data corresponding to entity
func (c *Component) Get(entity Entity) interface{} {
	if data, ok := c.data[entity]; ok {
		return data
	}
	return nil
}

// Not returns an inverted component used for filtering entities that don't have the component
func (c *Component) Not() *AntiComponent {
	return &AntiComponent{tag: c.tag}
}

func (c *Component) _Tag() *bit.Set {
	return &c.tag
}

func (c *Component) _Join(tag *bit.Set) *bit.Set {
	return tag.SetAnd(tag, &c.tag)
}

// AntiComponent is an inverted component used for filtering entities that don't have a component
type AntiComponent struct {
	tag bit.Set
}

func (a *AntiComponent) _Tag() *bit.Set {
	return &a.tag
}

func (a *AntiComponent) _Join(tag *bit.Set) *bit.Set {
	return tag.SetAndNot(tag, &a.tag)
}
