package rbtree

import (
	"golang.org/x/exp/constraints"
)

// ordered map like map[K]V
// implement by Tree
type Map[K, V any] struct {
	tree *rbTree[K, V]
}

// NewMap Create a new empty Map
func NewMap[K constraints.Ordered, V any]() Map[K, V] {
	cmp := func(a, b K) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	return Map[K, V]{tree: newTreeFunc[K, V](cmp)}
}

// Create a new empty map with compare function
func NewMapFunc[K, V any](cmp func(a, b K) int) Map[K, V] {
	return Map[K, V]{tree: newTreeFunc[K, V](cmp)}
}

// follow Map operation simple wrapper rbTree

func (m Map[K, V]) Len() int {
	return m.tree.Len()
}

func (m Map[K, V]) Min() MapIterator[K, V] {
	return newMapIterator(m.tree.Min())
}

func (m Map[K, V]) Max() MapIterator[K, V] {
	return newMapIterator(m.tree.Max())
}

func (m Map[K, V]) Limit() MapIterator[K, V] {
	return newMapIterator[K, V](nil)
}

func (m Map[K, V]) Find(k K) MapIterator[K, V] {
	return newMapIterator(m.tree.Find(k))
}

func (m Map[K, V]) FindGE(k K) MapIterator[K, V] {
	return MapIterator[K, V]{m.tree.FindGE(k)}
}

func (m Map[K, V]) FindLE(k K) MapIterator[K, V] {
	return MapIterator[K, V]{m.tree.FindLE(k)}
}

// Get from map
// value: value find with key
// ok: ture if key found
func (m Map[K, V]) Get(k K) (V, bool) {
	nd := m.tree.Find(k)
	if nd == nil {
		var ep V
		return ep, false
	}
	return nd.value, true
}

// Set key and value, create new pair if not exist
func (m Map[K, V]) Set(k K, v V) {
	nd, _ := m.tree.Insert(k)
	nd.value = v
}

// Delete an item with the given key.
// Return true if the item was found.
func (m Map[K, V]) DeleteWithKey(k K) bool {
	nd := m.tree.Find(k)
	if nd != nil {
		m.tree.Delete(nd)
		return true
	}
	return false
}

func (m Map[K, V]) DeleteWithIterator(iter MapIterator[K, V]) {
	m.tree.Delete(iter.nd)
}

// iterator allows scanning map elements in sort order.
//
// iterator invalidation rule is the same as C++ std::map<>'s. That
// is, if you delete the element that an iterator points to, the
// iterator becomes invalid. For other operation types, the iterator
// remains valid.
//
// default iterator is not valid, Limit() == false
type MapIterator[K, V any] struct {
	nd *node[K, V]
}

func newMapIterator[K, V any](n *node[K, V]) MapIterator[K, V] {
	return MapIterator[K, V]{n}
}

// allow clients to verify iterator is from the right map
// REQUIRES: !iter.Limit() && !iter.NegativeLimit()
func (iter MapIterator[K, V]) Map() Map[K, V] {
	return Map[K, V]{tree: iter.nd.myTree}
}

// Check if iterator equal, same as it1 == it2
func (iter MapIterator[K, V]) Equal(iter2 MapIterator[K, V]) bool {
	return iter.nd == iter2.nd
}

// Check if the iterator points to element in the map
func (iter MapIterator[K, V]) Limit() bool {
	return iter.nd == nil
}

// Create a new iterator that points to the successor of the current element.
// REQUIRES: !iter.Limit()
func (iter MapIterator[K, V]) Next() MapIterator[K, V] {
	return MapIterator[K, V]{iter.nd.doNext()}
}

// Create a new iterator that points to the predecessor of the current node.
// REQUIRES: !iter.Limit()
func (iter MapIterator[K, V]) Prev() MapIterator[K, V] {
	return MapIterator[K, V]{iter.nd.doPrev()}
}

// Return the current key in element.
// REQUIRES: !iter.Limit()
func (iter MapIterator[K, V]) Key() K {
	return iter.nd.item
}

// Return the current value in element.
// REQUIRES: !iter.Limit()
func (iter MapIterator[K, V]) Value() V {
	return iter.nd.value
}

// Return the current value pointer in element, nil if iter.Limit()
func (iter MapIterator[K, V]) ValuePointer() *V {
	if iter.nd == nil {
		return nil
	}
	return &iter.nd.value
}
