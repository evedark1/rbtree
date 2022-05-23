package rbtree

import (
	"fmt"
	"math/rand"
	"testing"
)

func testNewIntSet() Set[int] {
	return NewSet[int]()
}

func TestSetInsert(t *testing.T) {
	m := testNewIntSet()
	equalAssert(t, m.Len(), 0)
	equalAssert(t, true, m.Insert(0))
	equalAssert(t, true, m.Insert(1))
	equalAssert(t, true, m.Insert(2))
	equalAssert(t, false, m.Insert(0))
	equalAssert(t, m.Len(), 3)

	equalAssert(t, false, m.Exist(-1))
	equalAssert(t, true, m.Exist(0))
	equalAssert(t, true, m.Exist(1))
	equalAssert(t, true, m.Exist(2))

	equalAssert(t, m.DeleteWithKey(-1), false)
	equalAssert(t, m.DeleteWithKey(1), true)
	equalAssert(t, false, m.Exist(1))
}

func compareSetFull(t *testing.T, o *oracle, s Set[int]) {
	mi := s.Min()
	oi := o.FindGE(t, -1)
	for !mi.Limit() && !oi.Limit() {
		v := oi.Item()
		equalAssert(t, mi.Item(), v)

		mi = mi.Next()
		oi = oi.Next()
	}
	equalAssert(t, true, mi.Limit())
	equalAssert(t, true, oi.Limit())
}

func TestSetDelete(t *testing.T) {
	s := testNewIntSet()
	o := newOracle()
	for i := 0; i < 10; i++ {
		s.Insert(i)
		o.Insert(i)
	}
	compareSetFull(t, o, s)

	o.Delete(7)
	s.DeleteWithKey(7)
	compareSetFull(t, o, s)

	o.Delete(0)
	s.DeleteWithKey(0)
	compareSetFull(t, o, s)
}

func TestSetFind(t *testing.T) {
	m := testNewIntSet()
	m.Insert(0)
	m.Insert(2)
	m.Insert(3)
	m.Insert(7)
	m.Insert(9)

	iter := m.Find(3)
	equalAssert(t, iter.Item(), 3)
	equalAssert(t, m.Find(1), m.Limit())

	iter = m.FindGE(3)
	equalAssert(t, iter.Item(), 3)
	iter = m.FindGE(4)
	equalAssert(t, iter.Item(), 7)
	equalAssert(t, m.FindGE(10), m.Limit())

	iter = m.FindLE(3)
	equalAssert(t, iter.Item(), 3)
	iter = m.FindLE(4)
	equalAssert(t, iter.Item(), 3)
	iter = m.FindLE(10)
	equalAssert(t, iter.Item(), 9)
	equalAssert(t, m.FindLE(-1), m.Limit())
}

func TestSetOrder(t *testing.T) {
	// get random keys
	keys := make([]int, 0)
	for i := 0; i < 10; i++ {
		keys = append(keys, i)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	m := testNewIntSet()
	for _, v := range keys {
		m.Insert(v)
	}

	order := 0
	for iter := m.Min(); !iter.Limit(); iter = iter.Next() {
		equalAssert(t, order, iter.Item())
		order++
	}

	order--
	for iter := m.Max(); !iter.Limit(); iter = iter.Prev() {
		equalAssert(t, order, iter.Item())
		order--
	}
}

//
// Examples
//

func ExampleSetString() {
	set := NewSet[string]()
	set.Insert("value10")
	set.Insert("value12")

	ok := set.Exist("value10")
	fmt.Printf("Exist(10) -> {%t}\n", ok)
	ok = set.Exist("value11")
	fmt.Printf("Exist(11) -> {%t}\n", ok)

	// Find an element >= 11
	iter := set.FindGE("value11")
	fmt.Printf("FindGE(11) -> {%t %s}\n", iter.Limit(), iter.Item())

	// Find an element >= 13
	iter = set.FindGE("value13")
	if !iter.Limit() {
		panic("There should be no element >= 13")
	}

	// Output:
	// Exist(10) -> {true}
	// Exist(11) -> {false}
	// FindGE(11) -> {false value12}
}
