package ecs

import (
	"reflect"
	"testing"

	"github.com/yourbasic/bit"
)

type (
	data1 struct{ v int }
	data2 struct{ v int }
	data3 struct{ v int }
)

func setupNullComponent() (*Manager, DataComponent, DataComponent, DataComponent) {
	manager := NewManager()
	c1 := manager.NewNullComponent()
	c2 := manager.NewNullComponent()
	c3 := manager.NewNullComponent()
	return manager, c1, c2, c3
}

func setupSliceComponent() (*Manager, DataComponent, DataComponent, DataComponent) {
	manager := NewManager()
	c1 := manager.NewSliceComponent()
	c2 := manager.NewSliceComponent()
	c3 := manager.NewSliceComponent()
	return manager, c1, c2, c3
}

func setupDenseSliceComponent() (*Manager, DataComponent, DataComponent, DataComponent) {
	manager := NewManager()
	c1 := manager.NewDenseSliceComponent()
	c2 := manager.NewDenseSliceComponent()
	c3 := manager.NewDenseSliceComponent()
	return manager, c1, c2, c3
}

func setupMapComponent() (*Manager, DataComponent, DataComponent, DataComponent) {
	manager := NewManager()
	c1 := manager.NewMapComponent()
	c2 := manager.NewMapComponent()
	c3 := manager.NewMapComponent()
	return manager, c1, c2, c3
}

func setup(manager *Manager, c1, c2, c3 DataComponent) {
	manager.NewEntity()

	manager.NewEntity().AddComponent(c1, &data1{1})
	manager.NewEntity().AddComponent(c2, &data2{2})
	manager.NewEntity().AddComponent(c3, &data3{3})

	manager.NewEntity().AddComponent(c1, &data1{4}).AddComponent(c2, &data2{5})
	manager.NewEntity().AddComponent(c1, &data1{6}).AddComponent(c3, &data3{7})
	manager.NewEntity().AddComponent(c2, &data2{8}).AddComponent(c3, &data3{9})

	manager.NewEntity().AddComponent(c1, &data1{10}).AddComponent(c2, &data2{11}).AddComponent(c3, &data3{12})
}

func TestAll(t *testing.T) {
	var m *Manager
	var c1, c2, c3 DataComponent
	for _, fSetup := range []func() (*Manager, DataComponent, DataComponent, DataComponent){setupNullComponent, setupSliceComponent, setupDenseSliceComponent, setupMapComponent} {
		m, c1, c2, c3 = fSetup()
		setup(m, c1, c2, c3)
		testJoin(t, m, c1, c2, c3)

		m, c1, c2, c3 = fSetup()
		setup(m, c1, c2, c3)
		testEntity(t, m, c1, c2, c3)

		m, c1, c2, c3 = fSetup()
		setup(m, c1, c2, c3)
		testComponents(t, m, c1, c2, c3)

		m, c1, c2, c3 = fSetup()
		setup(m, c1, c2, c3)
		testUtils(t, m, c1, c2, c3)
	}
}

func testJoin(t *testing.T, m *Manager, c1, c2, c3 DataComponent) {
	for _, x := range []struct {
		test  string
		tag   *bit.Set
		value *bit.Set
	}{
		{"T0", m.Join(), bit.New(1, 2, 3, 4, 5, 6, 7)},

		{"T1", m.Join(c1), bit.New(1, 4, 5, 7)},
		{"T2", m.Join(c1, c1), bit.New(1, 4, 5, 7)},

		{"T3", m.Join(c2.Not()), bit.New(1, 3, 5)},
		{"T4", m.Join(c2.Not(), c2.Not()), bit.New(1, 3, 5)},

		{"T5", m.Join(c3, c3.Not()), bit.New()},
		{"T6", m.Join(c3.Not(), c3), bit.New()},

		{"T7", m.Join(c3, c2.Not(), c1), bit.New(5)},
		{"T8", m.Join(c1, c3, c2.Not()), bit.New(5)},
		{"T9", m.Join(c2.Not(), c1, c3), bit.New(5)},
	} {
		if !x.tag.Equal(x.value) {
			t.Errorf("Test %v: Join() has tag %v, wants %v", x.test, x.tag, x.value)
		}
	}
}

func testEntity(t *testing.T, m *Manager, c1, c2, c3 DataComponent) {
	if !m.entities.Equal(bit.New(1, 2, 3, 4, 5, 6, 7)) {
		t.Errorf("Wrong data %v, wants %v", m.entities, bit.New(1, 2, 3, 4, 5, 6, 7))
	}
	if s1 := m.Join(c1).Size(); s1 != 4 {
		t.Errorf("Wrong size %v, wants %v", s1, 4)
	}
	if s2 := m.Join(c2).Size(); s2 != 4 {
		t.Errorf("Wrong size %v, wants %v", s2, 4)
	}

	Entity(0).RemoveComponent(c3)
	Entity(3).RemoveComponent(c3)
	Entity(7).RemoveComponent(c3)
	if !m.entities.Equal(bit.New(1, 2, 4, 5, 6, 7)) {
		t.Errorf("Wrong data %v, wants %v", m.entities, bit.New(1, 2, 4, 5, 6, 7))
	}

	m.DeleteEntities(4, 1)

	if s1 := c1._Tag().Size(); s1 != 2 {
		t.Errorf("Wrong size %v, wants %v", s1, 2)
	}
	if s2 := c2._Tag().Size(); s2 != 3 {
		t.Errorf("Wrong size %v, wants %v", s2, 3)
	}
	if !m.entities.Equal(bit.New(2, 5, 6, 7)) {
		t.Errorf("Wrong data %v, wants %v", m.entities, bit.New(2, 5, 6, 7))
	}

	Entity(3).AddComponent(c3, 3)
	if !m.entities.Equal(bit.New(2, 3, 5, 6, 7)) {
		t.Errorf("Wrong data %v, wants %v", m.entities, bit.New(2, 3, 5, 6, 7))
	}

	m.DeleteAllEntities()

	for _, component := range m.components {
		if !component._Tag().Empty() {
			t.Errorf("Component is not empty")
		}
		if component._Manager() != m {
			t.Errorf("Manager is not set correctly")
		}
		if mc, ok := component.(*MapComponent); ok {
			if mc.data == nil {
				t.Errorf("Map data is nil")
			}
		}
	}

	if m.currentEntityIndex != 0 {
		t.Errorf("Current entity index is not 0")
	}
	if !m.entities.Empty() {
		t.Errorf("Entities cache is not empty")
	}
}

func testComponents(t *testing.T, m *Manager, c1, c2, c3 DataComponent) {
	_, null2 := c2.(*NullComponent)
	_, null3 := c3.(*NullComponent)

	if v := c1.Get(0); v != nil {
		t.Errorf("Wrong data %v, wants %v", v, nil)
	}

	if v := c1.Get(-1); v != nil {
		t.Errorf("Wrong data %v, wants %v", v, nil)
	}

	if v := c1.Get(99); v != nil {
		t.Errorf("Wrong data %v, wants %v", v, nil)
	}

	data4 := c2.Get(4)
	if !null2 {
		if data4 == nil {
			t.Errorf("Wrong data %v, wants %v", data4, 5)
		} else if data4.(*data2).v != 5 {
			t.Errorf("Wrong data %v, wants %v", data4.(*data2).v, 5)
		}
	} else if data4 != nil {
		t.Errorf("Wrong data %v, wants %v", data4, nil)
	}

	c3.Set(0, &data3{-1})
	if !Entity(0).HasComponent(c3) && c3.Get(0) != nil {
		t.Errorf("Set() must not set the value if the entity does not have the component")
	}

	c3.Set(7, &data3{-1})
	if !null3 {
		if v := c3.Get(7).(*data3).v; v != -1 {
			t.Errorf("Wrong data %v, wants %v", v, -1)
		}
	} else if c3.Get(7) != nil {
		t.Errorf("Wrong data %v, wants %v", c3.Get(7), nil)
	}

	c3.Set(-1, nil)
	c3.Set(99, nil)
}

func testUtils(t *testing.T, m *Manager, c1, c2, c3 DataComponent) {
	sum := 0
	m.Join(c1).Visit(Visit(func(entity Entity) { sum++ }))
	if sum != c1._Tag().Size() {
		t.Errorf("Wrong size %v, wants %v", sum, c1._Tag().Size())
	}

	if v := GetFirst(m.Join(c1, c1.Not())); v != nil {
		t.Errorf("Wrong data %v, wants %v", v, nil)
	}

	if v := GetFirst(m.Join(c3)); v == nil {
		t.Errorf("Wrong data %v, wants %v", v, 3)
	} else if *v != 3 {
		t.Errorf("Wrong data %v, wants %v", *v, 3)
	}
}

func TestDenseSliceComponent(t *testing.T) {
	manager := NewManager()
	c := manager.NewDenseSliceComponent()

	entities := make([]Entity, 12)
	for iEntity := range entities {
		entities[iEntity] = manager.NewEntity()
	}
	checkDenseSliceInvariants(t, c)

	for _, x := range []struct {
		test string
		f    func()
		res  DenseSliceComponent
	}{
		{"T0", func() { entities[11].AddComponent(c, -9) }, DenseSliceComponent{
			data:     []interface{}{-9},
			dataID:   []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0},
			entityID: []Entity{11},
		}},
		{"T1", func() { entities[11].RemoveComponent(c) }, DenseSliceComponent{
			data:     []interface{}{},
			dataID:   []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
			entityID: []Entity{},
		}},
		{"T2", func() { entities[7].AddComponent(c, 0) }, DenseSliceComponent{
			data:     []interface{}{0},
			dataID:   []int{-1, -1, -1, -1, -1, -1, -1, 0, -1, -1, -1, -1},
			entityID: []Entity{7},
		}},
		{"T3", func() { entities[5].AddComponent(c, 1) }, DenseSliceComponent{
			data:     []interface{}{0, 1},
			dataID:   []int{-1, -1, -1, -1, -1, 1, -1, 0, -1, -1, -1, -1},
			entityID: []Entity{7, 5},
		}},
		{"T4", func() { entities[0].AddComponent(c, 2) }, DenseSliceComponent{
			data:     []interface{}{0, 1, 2},
			dataID:   []int{2, -1, -1, -1, -1, 1, -1, 0, -1, -1, -1, -1},
			entityID: []Entity{7, 5, 0},
		}},
		{"T5", func() { entities[9].AddComponent(c, 33); c.Set(9, 3) }, DenseSliceComponent{
			data:     []interface{}{0, 1, 2, 3},
			dataID:   []int{2, -1, -1, -1, -1, 1, -1, 0, -1, 3, -1, -1},
			entityID: []Entity{7, 5, 0, 9},
		}},
		{"T6", func() { entities[3].AddComponent(c, 4) }, DenseSliceComponent{
			data:     []interface{}{0, 1, 2, 3, 4},
			dataID:   []int{2, -1, -1, 4, -1, 1, -1, 0, -1, 3, -1, -1},
			entityID: []Entity{7, 5, 0, 9, 3},
		}},
		{"T7", func() { entities[5].RemoveComponent(c) }, DenseSliceComponent{
			data:     []interface{}{0, 4, 2, 3},
			dataID:   []int{2, -1, -1, 1, -1, -1, -1, 0, -1, 3, -1, -1},
			entityID: []Entity{7, 3, 0, 9},
		}},
		{"T8", func() { entities[5].AddComponent(c, 5) }, DenseSliceComponent{
			data:     []interface{}{0, 4, 2, 3, 5},
			dataID:   []int{2, -1, -1, 1, -1, 4, -1, 0, -1, 3, -1, -1},
			entityID: []Entity{7, 3, 0, 9, 5},
		}},
	} {
		x.f()
		checkDenseSliceInvariants(t, c)
		if !reflect.DeepEqual(c.data, x.res.data) {
			t.Errorf("Test %v: Wrong data %v, wants %v", x.test, c.data, x.res.data)
		}
		if !reflect.DeepEqual(c.dataID, x.res.dataID) {
			t.Errorf("Test %v: Wrong data %v, wants %v", x.test, c.dataID, x.res.dataID)
		}
		if !reflect.DeepEqual(c.entityID, x.res.entityID) {
			t.Errorf("Test %v: Wrong data %v, wants %v", x.test, c.entityID, x.res.entityID)
		}
	}
}

func checkDenseSliceInvariants(t *testing.T, c *DenseSliceComponent) {
	for iEntity := range c.entityID {
		if c.dataID[c.entityID[iEntity]] != iEntity {
			t.Errorf("Wrong data %v, wants %v", c.dataID[c.entityID[iEntity]], iEntity)
		}
	}
	for iData := range c.dataID {
		if c.dataID[iData] != -1 && int(c.entityID[c.dataID[iData]]) != iData {
			t.Errorf("Wrong data %v, wants %v", c.entityID[c.dataID[iData]], iData)
		}
	}
}

func TestMaintain(t *testing.T) {
	var manager *Manager
	var c1, c2, c3 DataComponent
	for _, fSetup := range []func() (*Manager, DataComponent, DataComponent, DataComponent){setupSliceComponent, setupDenseSliceComponent, setupMapComponent} {
		manager, c1, c2, c3 = fSetup()
		manager.Maintain(0, 1)

		manager.NewEntity()
		manager.NewEntity()

		manager.NewEntity().AddComponent(c1, &data1{1})
		manager.NewEntity()
		manager.NewEntity()
		manager.NewEntity().AddComponent(c2, &data2{2})
		manager.NewEntity()
		manager.NewEntity()
		manager.NewEntity().AddComponent(c3, &data3{3})
		manager.NewEntity()
		manager.NewEntity()

		manager.NewEntity().AddComponent(c1, &data1{4}).AddComponent(c2, &data2{5})
		manager.NewEntity()
		manager.NewEntity()
		manager.NewEntity().AddComponent(c1, &data1{6}).AddComponent(c3, &data3{7})
		manager.NewEntity()
		manager.NewEntity()
		manager.NewEntity().AddComponent(c2, &data2{8}).AddComponent(c3, &data3{9})
		manager.NewEntity()
		manager.NewEntity()

		manager.NewEntity().AddComponent(c1, &data1{10}).AddComponent(c2, &data2{11}).AddComponent(c3, &data3{12})
		manager.NewEntity()
		manager.NewEntity()

		manager.Maintain(0, 1)

		manager.NewEntity().AddComponent(c1, &data1{13})

		for _, x := range []struct {
			test  string
			data  int
			value int
		}{
			{"T00", c1.Get(0).(*data1).v, 10},
			{"T01", c2.Get(0).(*data2).v, 11},
			{"T02", c3.Get(0).(*data3).v, 12},

			{"T03", c2.Get(1).(*data2).v, 8},
			{"T04", c3.Get(1).(*data3).v, 9},

			{"T05", c1.Get(2).(*data1).v, 1},

			{"T06", c1.Get(3).(*data1).v, 6},
			{"T07", c3.Get(3).(*data3).v, 7},

			{"T08", c1.Get(4).(*data1).v, 4},
			{"T09", c2.Get(4).(*data2).v, 5},

			{"T10", c2.Get(5).(*data2).v, 2},

			{"T11", c3.Get(6).(*data3).v, 3},

			{"T12", c1.Get(7).(*data1).v, 13},
		} {
			if x.data != x.value {
				t.Errorf("Test %v: Wrong data %v, wants %v", x.test, x.data, x.value)
			}
		}

		manager.Maintain(0, 0)
	}
}
