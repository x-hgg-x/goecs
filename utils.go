package ecs

// Visit is a decorator function for bit.Set.Visit() method
func Visit(f func(entity Entity)) func(index int) bool {
	return func(index int) bool {
		f(Entity(index))
		return false
	}
}
