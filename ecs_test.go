package ecs

import (
	"testing"

	"github.com/yourbasic/bit"
)

type (
	c1 struct{ v int }
	c2 struct{ v int }
	c3 struct{ v int }
)

type testComponents struct {
	c1 *Component
	c2 *Component
	c3 *Component
}

func setup() (Manager, testComponents) {
	manager := Manager{}

	components := testComponents{
		c1: manager.NewComponent(),
		c2: manager.NewComponent(),
		c3: manager.NewComponent(),
	}

	manager.NewEntity()

	manager.NewEntity().AddComponent(components.c1, &c1{1})
	manager.NewEntity().AddComponent(components.c2, &c2{2})
	manager.NewEntity().AddComponent(components.c3, &c3{3})

	manager.NewEntity().AddComponent(components.c1, &c1{4}).AddComponent(components.c2, &c2{5})
	manager.NewEntity().AddComponent(components.c1, &c1{6}).AddComponent(components.c3, &c3{7})
	manager.NewEntity().AddComponent(components.c2, &c2{8}).AddComponent(components.c3, &c3{9})

	manager.NewEntity().AddComponent(components.c1, &c1{10}).AddComponent(components.c2, &c2{11}).AddComponent(components.c3, &c3{12})

	return manager, components
}

func TestJoin(t *testing.T) {
	m, c := setup()

	for _, x := range []struct {
		test  string
		tag   *bit.Set
		value *bit.Set
	}{
		{"T0", m.Join(), bit.New(1, 2, 3, 4, 5, 6, 7)},

		{"T1", m.Join(c.c1), bit.New(1, 4, 5, 7)},
		{"T2", m.Join(c.c1, c.c1), bit.New(1, 4, 5, 7)},

		{"T3", m.Join(c.c2.Not()), bit.New(1, 3, 5)},
		{"T4", m.Join(c.c2.Not(), c.c2.Not()), bit.New(1, 3, 5)},

		{"T5", m.Join(c.c3, c.c3.Not()), bit.New()},
		{"T6", m.Join(c.c3.Not(), c.c3), bit.New()},

		{"T7", m.Join(c.c3, c.c2.Not(), c.c1), bit.New(5)},
		{"T8", m.Join(c.c1, c.c3, c.c2.Not()), bit.New(5)},
		{"T9", m.Join(c.c2.Not(), c.c1, c.c3), bit.New(5)},
	} {
		if !x.tag.Equal(x.value) {
			t.Errorf("Test %v: Join() has tag %v, wants %v", x.test, x.tag, x.value)
		}
	}
}

func TestEntity(t *testing.T) {
	m, c := setup()

	if s1 := m.Join(c.c1).Size(); s1 != 4 {
		t.Errorf("Wrong size %v, wants %v", s1, 4)
	}
	if s2 := m.Join(c.c2).Size(); s2 != 4 {
		t.Errorf("Wrong size %v, wants %v", s2, 4)
	}

	m.DeleteEntities(4, 1)

	if s1 := c.c1.tag.Size(); s1 != 2 {
		t.Errorf("Wrong size %v, wants %v", s1, 2)
	}
	if s2 := c.c2.tag.Size(); s2 != 3 {
		t.Errorf("Wrong size %v, wants %v", s2, 3)
	}

	m.DeleteAllEntities()

	for _, component := range m.components {
		if !component.tag.Empty() {
			t.Errorf("Component is not empty")
		}
		if component.data == nil || len(component.data) != 0 {
			t.Errorf("Data is nil or not empty")
		}
	}

	if m.currentEntityIndex != 0 {
		t.Errorf("Current entity index is not 0")
	}
}

func TestComponents(t *testing.T) {
	_, c := setup()

	if !c.c1._Tag().Equal(&c.c1.tag) || !c.c1.Not()._Tag().Equal(&c.c1.tag) {
		t.Errorf("_Tag() method is incorrect")
	}

	if v := c.c1.Get(0); v != nil {
		t.Errorf("Wrong data %v, wants %v", v, nil)
	}

	if v := c.c2.Get(4); v == nil {
		t.Errorf("Wrong data %v, wants %v", v, 5)
	} else if entity := v.(*c2); entity.v != 5 {
		t.Errorf("Wrong data %v, wants %v", entity.v, 5)
	}

	c.c3.Set(0, &c3{-1})
	if Entity(0).HasComponent(c.c3) {
		t.Errorf("Set() must not set the value if the entity does not have the component")
	}

	c.c3.Set(7, &c3{-1})
	if v := c.c3.Get(7).(*c3).v; v != -1 {
		t.Errorf("Wrong data %v, wants %v", v, -1)
	}
}

func TestUtils(t *testing.T) {
	m, c := setup()

	sum := 0
	m.Join(c.c1).Visit(Visit(func(entity Entity) { sum++ }))
	if sum != c.c1.tag.Size() {
		t.Errorf("Wrong size %v, wants %v", sum, c.c1.tag.Size())
	}

	if v := GetFirst(m.Join(c.c1, c.c1.Not())); v != nil {
		t.Errorf("Wrong data %v, wants %v", v, nil)
	}

	if v := GetFirst(m.Join(c.c3)); v == nil {
		t.Errorf("Wrong data %v, wants %v", v, 3)
	} else if *v != 3 {
		t.Errorf("Wrong data %v, wants %v", *v, 3)
	}
}
