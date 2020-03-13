# GoEcs

An implementation of the ECS paradigm in Go.

## Usage

```go
package main

import (
	"fmt"

	ecs "github.com/x-hgg-x/goecs"
)

func main() {
	// List of component data types
	type Shape struct{ shape string }
	type Color struct{ color string }
	type Name struct{ name string }

	// Structure for storing components
	components := struct {
		Shape *ecs.Component
		Color *ecs.Component
		Name  *ecs.Component
		Value *ecs.Component
	}{}

	// Initialize a new manager
	manager := ecs.NewManager()

	// Create components
	components.Shape = manager.NewComponent()
	components.Color = manager.NewComponent()
	components.Name = manager.NewComponent()
	components.Value = manager.NewComponent()

	// Create entities
	manager.NewEntity().AddComponent(components.Shape, &Shape{"square"}).AddComponent(components.Color, &Color{"red"})
	manager.NewEntity().AddComponent(components.Shape, &Shape{"circle"}).AddComponent(components.Name, &Name{"tom"})
	manager.NewEntity().AddComponent(components.Color, &Color{"blue"}).AddComponent(components.Name, &Name{"john"})

	manager.NewEntity().
		AddComponent(components.Shape, &Shape{"triangle"}).
		AddComponent(components.Color, &Color{"green"}).
		AddComponent(components.Name, &Name{"paul"})

	// Loop on entities which have specified components
	// The Join() method gives a bit.Set tag containing integers which can be converted to entities,
	// and we use the bit.Set.Visit() method to loop through the set.
	// The decorator ecs.Visit() is used when we want to iterate through all elements of the set.
	// It converts each set element to an entity.
	manager.Join(components.Shape, components.Name).Visit(ecs.Visit(func(entity ecs.Entity) {
		shape := components.Shape.Get(entity).(*Shape)
		name := components.Name.Get(entity).(*Name)
		fmt.Printf("Entity has the shape '%s' and the name '%s'\n", shape.shape, name.name)
	}))
	fmt.Println()

	// If we want to break the loop when some condition is met, we use the Visit() method directly
	aborted := manager.Join(components.Shape).Visit(func(index int) (skip bool) {
		shape := components.Shape.Get(ecs.Entity(index)).(*Shape)
		fmt.Printf("Entity has the shape '%s'\n", shape.shape)
		if shape.shape == "circle" {
			shape.shape = "CIRCLE"
			fmt.Printf("Entity has now the shape '%s'\n", shape.shape)
			return true
		}
		return false
	})
	fmt.Printf("Loop aborted: %v\n\n", aborted)

	// The helper function ecs.GetFirst() is useful if we want only the first entity matching a tag
	if firstEntity := ecs.GetFirst(manager.Join(components.Shape, components.Color, components.Name)); firstEntity != nil {
		shape := components.Shape.Get(*firstEntity).(*Shape)
		color := components.Color.Get(*firstEntity).(*Color)
		name := components.Name.Get(*firstEntity).(*Name)
		fmt.Printf("First matching entity has the shape '%s', the color '%s' and the name '%s'\n", shape.shape, color.color, name.name)
	}
	fmt.Println()

	// The Not() method is used when we want to exclude a particular component
	manager.Join(components.Shape.Not()).Visit(ecs.Visit(func(entity ecs.Entity) {
		fmt.Printf("Entity components: ")
		if entity.HasComponent(components.Color) {
			fmt.Printf("Color: '%s', ", components.Color.Get(entity).(*Color).color)
		}
		if entity.HasComponent(components.Name) {
			fmt.Printf("Name: '%s'", components.Name.Get(entity).(*Name).name)
		}
		fmt.Println()
	}))
	fmt.Println()

	// To iterate through all entities with at least one component, we use the Join() method without any argument
	manager.Join().Visit(ecs.Visit(func(entity ecs.Entity) {
		fmt.Printf("Entity components: ")
		if entity.HasComponent(components.Shape) {
			fmt.Printf("Shape: '%s', ", components.Shape.Get(entity).(*Shape).shape)
		}
		if entity.HasComponent(components.Color) {
			fmt.Printf("Color: '%s', ", components.Color.Get(entity).(*Color).color)
		}
		if entity.HasComponent(components.Name) {
			fmt.Printf("Name: '%s'", components.Name.Get(entity).(*Name).name)
		}
		fmt.Println()
	}))
	fmt.Println()

	// If the component data is not a pointer, we can use the Set() method to change its value
	manager.NewEntity().AddComponent(components.Value, 3)
	firstEntity := *ecs.GetFirst(manager.Join(components.Value))
	fmt.Println("Old value:", components.Value.Get(firstEntity).(int))
	components.Value.Set(firstEntity, 4)
	fmt.Println("New value:", components.Value.Get(firstEntity).(int))
}
```
