package rbtree

import (
	"fmt"
	"math/rand"
	"testing"
)

func testNewIntMap() Map[int, int] {
	return NewMap[int, int]()
}

func TestMapSetGet(t *testing.T) {
	m := testNewIntMap()
	equalAssert(t, m.Len(), 0)
	m.Set(0, 10)
	m.Set(1, 11)
	m.Set(2, 12)
	m.Set(0, 13)
	equalAssert(t, m.Len(), 3)

	var val int
	var ok bool

	val, ok = m.Get(-1)
	equalAssert(t, val, 0)
	equalAssert(t, ok, false)

	val, ok = m.Get(0)
	equalAssert(t, val, 13)
	equalAssert(t, ok, true)

	val, ok = m.Get(1)
	equalAssert(t, val, 11)
	equalAssert(t, ok, true)

	val, ok = m.Get(2)
	equalAssert(t, val, 12)
	equalAssert(t, ok, true)

	equalAssert(t, m.DeleteWithKey(-1), false)
	equalAssert(t, m.DeleteWithKey(1), true)
	_, ok = m.Get(1)
	equalAssert(t, ok, false)
}

func compareMapFull(t *testing.T, o *oracle, m Map[int, int]) {
	mi := m.Min()
	oi := o.FindGE(t, -1)
	for !mi.Limit() && !oi.Limit() {
		v := oi.Item()
		equalAssert(t, mi.Key(), v)
		equalAssert(t, mi.Value(), v)

		mi = mi.Next()
		oi = oi.Next()
	}
	equalAssert(t, true, mi.Limit())
	equalAssert(t, true, oi.Limit())
}

func TestMapDelete(t *testing.T) {
	m := testNewIntMap()
	o := newOracle()
	for i := 0; i < 10; i++ {
		m.Set(i, i)
		o.Insert(i)
	}
	compareMapFull(t, o, m)

	o.Delete(7)
	m.DeleteWithKey(7)
	compareMapFull(t, o, m)

	o.Delete(0)
	m.DeleteWithKey(0)
	compareMapFull(t, o, m)
}

func TestMapFind(t *testing.T) {
	m := testNewIntMap()
	m.Set(0, 0)
	m.Set(2, 1)
	m.Set(3, 2)
	m.Set(7, 3)
	m.Set(9, 4)

	iter := m.Find(3)
	equalAssert(t, iter.Key(), 3)
	equalAssert(t, m.Find(1), m.Limit())

	iter = m.FindGE(3)
	equalAssert(t, iter.Key(), 3)
	iter = m.FindGE(4)
	equalAssert(t, iter.Key(), 7)
	equalAssert(t, m.FindGE(10), m.Limit())

	iter = m.FindLE(3)
	equalAssert(t, iter.Key(), 3)
	iter = m.FindLE(4)
	equalAssert(t, iter.Key(), 3)
	iter = m.FindLE(10)
	equalAssert(t, iter.Key(), 9)
	equalAssert(t, m.FindLE(-1), m.Limit())
}

func TestMapOrder(t *testing.T) {
	// get random keys
	keys := make([]int, 0)
	for i := 0; i < 10; i++ {
		keys = append(keys, i)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	m := testNewIntMap()
	for _, v := range keys {
		m.Set(v, v)
	}

	order := 0
	for iter := m.Min(); !iter.Limit(); iter = iter.Next() {
		equalAssert(t, order, iter.Key())
		equalAssert(t, order, iter.Value())
		order++
	}

	order--
	for iter := m.Max(); !iter.Limit(); iter = iter.Prev() {
		equalAssert(t, order, iter.Key())
		equalAssert(t, order, iter.Value())
		order--
	}
}

//
// Examples
//

func ExampleMapIntString() {
	tree := NewMap[int, string]()
	tree.Set(10, "value10")
	tree.Set(12, "value12")

	val, ok := tree.Get(10)
	fmt.Printf("Get(10) -> {%t %s}\n", ok, val)
	val, ok = tree.Get(11)
	fmt.Printf("Get(11) -> {%t %s}\n", ok, val)

	// Find an element >= 11
	iter := tree.FindGE(11)
	fmt.Printf("FindGE(11) -> {%t %s}\n", iter.Limit(), iter.Value())

	// Find an element >= 13
	iter = tree.FindGE(13)
	if !iter.Limit() {
		panic("There should be no element >= 13")
	}

	// Output:
	// Get(10) -> {true value10}
	// Get(11) -> {false }
	// FindGE(11) -> {false value12}
}
