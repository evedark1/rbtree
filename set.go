package rbtree

import (
	"golang.org/x/exp/constraints"
)

// ordered set like std::set<K>
// implement by Tree
type Set[K any] struct {
	tree *rbTree[K, voidType]
}

// NewSet Create a new empty Set
func NewSet[K constraints.Ordered]() Set[K] {
	cmp := func(a, b K) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	return Set[K]{tree: newTreeFunc[K, voidType](cmp)}
}

// Create a new empty set with compare function
func NewSetFunc[K any](cmp func(a, b K) int) Set[K] {
	return Set[K]{tree: newTreeFunc[K, voidType](cmp)}
}

// follow Map operation simple wrapper rbTree

func (m Set[K]) Len() int {
	return m.tree.Len()
}

func (m Set[K]) Min() SetIterator[K] {
	return newSetIterator(m.tree.Min())
}

func (m Set[K]) Max() SetIterator[K] {
	return newSetIterator(m.tree.Max())
}

func (m Set[K]) Limit() SetIterator[K] {
	return newSetIterator[K](nil)
}

func (m Set[K]) Find(k K) SetIterator[K] {
	return newSetIterator(m.tree.Find(k))
}

func (m Set[K]) FindGE(k K) SetIterator[K] {
	return SetIterator[K]{m.tree.FindGE(k)}
}

func (m Set[K]) FindLE(k K) SetIterator[K] {
	return SetIterator[K]{m.tree.FindLE(k)}
}

// Get from map
// value: value find with key
// ok: ture if key found
func (m Set[K]) Exist(k K) bool {
	return m.tree.Find(k) != nil
}

// Insert key into set
// return true if key is inserted or false if key is existed
func (m Set[K]) Insert(k K) bool {
	_, exist := m.tree.Insert(k)
	return exist
}

// Delete an item with the given key.
// Return true if the item was found.
func (m Set[K]) DeleteWithKey(k K) bool {
	nd := m.tree.Find(k)
	if nd != nil {
		m.tree.Delete(nd)
		return true
	}
	return false
}

func (m Set[K]) DeleteWithIterator(iter SetIterator[K]) {
	m.tree.Delete(iter.nd)
}

// iterator allows scanning set elements in sort order.
//
// iterator invalidation rule is the same as C++ std::set<>'s. That
// is, if you delete the element that an iterator points to, the
// iterator becomes invalid. For other operation types, the iterator
// remains valid.
//
// default iterator is not valid, Limit() == false
type SetIterator[K any] struct {
	nd *node[K, voidType]
}

func newSetIterator[K any](n *node[K, voidType]) SetIterator[K] {
	return SetIterator[K]{n}
}

// allow clients to verify iterator is from the right map
// REQUIRES: !iter.Limit() && !iter.NegativeLimit()
func (iter SetIterator[K]) Set() Set[K] {
	return Set[K]{tree: iter.nd.myTree}
}

// Check if iterator equal, same as it1 == it2
func (iter SetIterator[K]) Equal(iter2 SetIterator[K]) bool {
	return iter.nd == iter2.nd
}

// Check if the iterator points to element in the map
func (iter SetIterator[K]) Limit() bool {
	return iter.nd == nil
}

// Create a new iterator that points to the successor of the current element.
// REQUIRES: !iter.Limit()
func (iter SetIterator[K]) Next() SetIterator[K] {
	return SetIterator[K]{iter.nd.doNext()}
}

// Create a new iterator that points to the predecessor of the current node.
// REQUIRES: !iter.Limit()
func (iter SetIterator[K]) Prev() SetIterator[K] {
	return SetIterator[K]{iter.nd.doPrev()}
}

// Return the current item in element.
// REQUIRES: !iter.Limit()
func (iter SetIterator[K]) Item() K {
	return iter.nd.item
}
