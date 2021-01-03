package ecs

import "testing"

func BenchmarkSliceComponent(b *testing.B) {
	manager := NewManager()
	c := manager.NewSliceComponent()

	for i := 0; i < 10000; i++ {
		manager.NewEntity().AddComponent(c, 1.0)
	}

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Get(Entity(i % 10000))
		}
	})

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Set(Entity(i%10000), 2.0)
		}
	})

	b.Run("Remove", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c._Remove(Entity(i % 10000))
		}
	})
}

func BenchmarkDenseSliceComponent(b *testing.B) {
	manager := NewManager()

	c := manager.NewDenseSliceComponent()
	for i := 0; i < 10000; i++ {
		manager.NewEntity().AddComponent(c, 1.0)
	}

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Get(Entity(i % 10000))
		}
	})

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Set(Entity(i%10000), 2.0)
		}
	})

	b.Run("Remove", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c._Remove(Entity(i % 10000))
		}
	})
}

func BenchmarkMapComponent(b *testing.B) {
	manager := NewManager()

	c := manager.NewMapComponent()
	for i := 0; i < 10000; i++ {
		manager.NewEntity().AddComponent(c, 1.0)
	}

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Get(Entity(i % 10000))
		}
	})

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Set(Entity(i%10000), 2.0)
		}
	})

	b.Run("Remove", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c._Remove(Entity(0))
		}
	})
}
