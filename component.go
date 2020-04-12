package ecs

import "github.com/yourbasic/bit"

// DataComponent is a component with data
type DataComponent interface {
	joinable
	_Tag() *bit.Set
	_Manager() *Manager
	_Reset()
	_Remove(Entity)
	Get(Entity) interface{}
	Set(Entity, interface{})
	Not() *AntiComponent
}

type joinable interface {
	_Join(*bit.Set) *bit.Set
}

type component struct {
	tag     bit.Set
	manager *Manager
}

func (c *component) _Tag() *bit.Set {
	return &c.tag
}

func (c *component) _Manager() *Manager {
	return c.manager
}

func (c *component) _Join(tag *bit.Set) *bit.Set {
	return tag.SetAnd(tag, &c.tag)
}

// Not returns an inverted component used for filtering entities that don't have the component
func (c *component) Not() *AntiComponent {
	return &AntiComponent{tag: c.tag}
}

// AntiComponent is an inverted component used for filtering entities that don't have a component
type AntiComponent struct {
	tag bit.Set
}

func (a *AntiComponent) _Join(tag *bit.Set) *bit.Set {
	return tag.SetAndNot(tag, &a.tag)
}
