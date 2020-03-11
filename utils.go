package ecs

import "github.com/yourbasic/bit"

// Visit is a decorator function for bit.Set.Visit() method
func Visit(f func(entity Entity)) func(index int) bool {
	return func(index int) bool {
		f(Entity(index))
		return false
	}
}

// GetFirst returns a reference to the first entity matching a tag or nil if there are none
func GetFirst(tag *bit.Set) *Entity {
	if tag.Empty() {
		return nil
	}

	firstEntity := Entity(tag.Next(-1))
	return &firstEntity
}
