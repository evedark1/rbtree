//
// Created by Yaz Saito on 06/10/12.
//

// A red-black tree with an API similar to C++ STL's.
//
// The implementation is inspired (read: stolen) from:
// http://en.literateprograms.org/Red-black_tree_(C)#chunk use:private function prototypes.
//
package rbtree

import (
	"fmt"
	"strings"
)

type voidType struct{}

type rbTree[K, V any] struct {
	// Root of the tree
	root *node[K, V]

	// The minimum and maximum nodes under the root.
	minNode, maxNode *node[K, V]

	// Number of nodes under root, including the root
	count int

	// CompareFunc returns 0 if a==b, <0 if a<b, >0 if a>b.
	compare func(a, b K) int
}

// Create a new empty tree with compare function
func newTreeFunc[K, V any](cmp func(a, b K) int) *rbTree[K, V] {
	return &rbTree[K, V]{compare: cmp}
}

// Return the number of elements in the tree.
func (root *rbTree[K, V]) Len() int {
	return root.count
}

// find the minimum node in the tree
// If the tree is empty, return nil
func (root *rbTree[K, V]) Min() *node[K, V] {
	return root.minNode
}

// find the maximum node in the tree
// If the tree is empty, return nil
func (root *rbTree[K, V]) Max() *node[K, V] {
	return root.maxNode
}

// A convenience function for finding an node equal to key, and return nil if not found.
// DO NOT modify return node except node.value
func (root *rbTree[K, V]) Find(key K) *node[K, V] {
	n, exact := root.findGE(key)
	if exact {
		return n
	}
	return nil
}

// Find the smallest element N such that N >= key, and return the pointer to the node.
// If no such element is found, return nil
func (root *rbTree[K, V]) FindGE(key K) *node[K, V] {
	n, _ := root.findGE(key)
	return n
}

// Find the largest element N such that N <= key, and return the pointer to the node.
// If no such element is found, return nil
func (root *rbTree[K, V]) FindLE(key K) *node[K, V] {
	n, exact := root.findGE(key)
	if exact {
		return n
	}
	if n != nil {
		return n.doPrev()
	}
	return root.maxNode
}

func getGU[K, V any](n *node[K, V]) (grandparent, uncle *node[K, V]) {
	grandparent = n.parent.parent
	if n.parent.isLeftChild() {
		uncle = grandparent.right
	} else {
		uncle = grandparent.left
	}
	return
}

// Insert an item.
// if the item is already in the tree, return node, false
// Otherwise add a new node, return node, true
func (root *rbTree[K, V]) Insert(key K) (*node[K, V], bool) {

	rn, ok := root.doInsert(key)
	if !ok {
		return rn, false
	}

	n := rn
	n.color = red
	var uncle, grandparent *node[K, V]
	for {

		// Case 1: N is at the root
		if n.parent == nil {
			n.color = black
			break
		}

		// Case 2: The parent is black, so the tree already
		// satisfies the RB properties
		if n.parent.color == black {
			break
		}

		// Case 3: parent and uncle are both red.
		// Then paint both black and make grandparent red.
		grandparent, uncle = getGU(n)

		if uncle != nil && uncle.color == red {
			n.parent.color = black
			uncle.color = black
			grandparent.color = red
			n = grandparent
			continue
		}

		// Case 4: parent is red, uncle is black (1)
		if n.isRightChild() && n.parent.isLeftChild() {
			root.rotateLeft(n.parent)

			n = n.left
			grandparent, uncle = getGU(n)
			//continue
		} else {
			if n.isLeftChild() && n.parent.isRightChild() {
				root.rotateRight(n.parent)
				n = n.right
				grandparent, uncle = getGU(n)
				//continue
			}
		}

		// Case 5: parent is red, uncle is black (2)
		n.parent.color = black
		grandparent.color = red

		if n.isLeftChild() && n.parent.isLeftChild() {
			root.rotateRight(grandparent)
		} else {
			if n.isRightChild() && n.parent.isRightChild() {
				root.rotateLeft(grandparent)
			} else {
				panic(fmt.Sprintf("assertion fails: should not get here on case 5."))
			}
		}
		break
	}
	return rn, true
}

// Delete an item with the given key. Return true iff the item was
// found.
func (root *rbTree[K, V]) DeleteWithKey(key K) bool {
	n, exact := root.findGE(key)
	if exact {
		root.doDelete(n)
		return true
	}
	return false
}

// Delete the current node.
func (root *rbTree[K, V]) Delete(node *node[K, V]) {
	doAssert(node != nil && node.myTree == root)
	root.doDelete(node)
}

func doAssert(b bool) {
	if !b {
		panic("rbtree internal assertion failed")
	}
}

const red = iota
const black = 1 + iota

type node[K, V any] struct {
	myTree              *rbTree[K, V]
	parent, left, right *node[K, V]
	color               int // black or red
	item                K
	value               V
}

//
// Internal node attribute accessors
//
func getColor[K, V any](n *node[K, V]) int {
	if n == nil {
		return black
	}
	return n.color
}

func (n *node[K, V]) isLeftChild() bool {
	return n == n.parent.left
}

func (n *node[K, V]) isRightChild() bool {
	return n == n.parent.right
}

func (n *node[K, V]) sibling() *node[K, V] {
	doAssert(n.parent != nil)
	if n.isLeftChild() {
		return n.parent.right
	}
	return n.parent.left
}

// Return the minimum node that's larger than N. Return nil if no such
// node is found.
func (n *node[K, V]) doNext() *node[K, V] {
	if n.right != nil {
		m := n.right
		for m.left != nil {
			m = m.left
		}
		return m
	}

	for n != nil {
		p := n.parent
		if p == nil {
			return nil
		}
		if n.isLeftChild() {
			return p
		}
		n = p
	}
	return nil
}

// Return the maximum node that's smaller than N. Return nil if no
// such node is found.
func (n *node[K, V]) doPrev() *node[K, V] {
	if n.left != nil {
		return maxPredecessor(n)
	}

	for n != nil {
		p := n.parent
		if p == nil {
			break
		}
		if n.isRightChild() {
			return p
		}
		n = p
	}
	return nil
}

// Return the predecessor of "n".
func maxPredecessor[K, V any](n *node[K, V]) *node[K, V] {
	doAssert(n.left != nil)
	m := n.left
	for m.right != nil {
		m = m.right
	}
	return m
}

//
// rbTree methods
//

//
// Private methods
//

func (root *rbTree[K, V]) recomputeMinNode() {
	root.minNode = root.root
	if root.minNode != nil {
		for root.minNode.left != nil {
			root.minNode = root.minNode.left
		}
	}
}

func (root *rbTree[K, V]) recomputeMaxNode() {
	root.maxNode = root.root
	if root.maxNode != nil {
		for root.maxNode.right != nil {
			root.maxNode = root.maxNode.right
		}
	}
}

func (root *rbTree[K, V]) maybeSetMinNode(n *node[K, V]) {
	if root.minNode == nil {
		root.minNode = n
		root.maxNode = n
	} else if root.compare(n.item, root.minNode.item) < 0 {
		root.minNode = n
	}
}

func (root *rbTree[K, V]) maybeSetMaxNode(n *node[K, V]) {
	if root.maxNode == nil {
		root.minNode = n
		root.maxNode = n
	} else if root.compare(n.item, root.maxNode.item) > 0 {
		root.maxNode = n
	}
}

// Try inserting "item" into the tree.
// if the item is already in the tree, return node, false
// Otherwise add a new (leaf) node, return node, true
func (root *rbTree[K, V]) doInsert(item K) (*node[K, V], bool) {
	if root.root == nil {
		n := &node[K, V]{item: item, myTree: root}
		root.root = n
		root.minNode = n
		root.maxNode = n
		root.count++
		return n, true
	}
	parent := root.root
	for true {
		comp := root.compare(item, parent.item)
		if comp == 0 {
			return parent, false
		} else if comp < 0 {
			if parent.left == nil {
				n := &node[K, V]{item: item, parent: parent, myTree: root}
				parent.left = n
				root.count++
				root.maybeSetMinNode(n)
				return n, true
			} else {
				parent = parent.left
			}
		} else {
			if parent.right == nil {
				n := &node[K, V]{item: item, parent: parent, myTree: root}
				parent.right = n
				root.count++
				root.maybeSetMaxNode(n)
				return n, true
			} else {
				parent = parent.right
			}
		}
	}
	panic("should not reach here")
}

// Find a node whose item >= key. The 2nd return value is true iff the
// node.item==key. Returns (nil, false) if all nodes in the tree are <
// key.
func (root *rbTree[K, V]) findGE(key K) (*node[K, V], bool) {
	n := root.root
	for true {
		if n == nil {
			return nil, false
		}
		comp := root.compare(key, n.item)
		if comp == 0 {
			return n, true
		} else if comp < 0 {
			if n.left != nil {
				n = n.left
			} else {
				return n, false
			}
		} else {
			if n.right != nil {
				n = n.right
			} else {
				succ := n.doNext()
				if succ == nil {
					return nil, false
				} else {
					comp = root.compare(key, succ.item)
					return succ, (comp == 0)
				}
			}
		}
	}
	panic("should not reach here")
}

// Delete N from the tree.
func (root *rbTree[K, V]) doDelete(n *node[K, V]) {
	if n.myTree != nil && n.myTree != root {
		panic(fmt.Sprintf("delete applied to node that was not from our tree... n has tree: '%s'\n\n while root has tree: '%s'\n\n", n.myTree.DumpAsString(), root.DumpAsString()))
	}
	if n.left != nil && n.right != nil {
		pred := maxPredecessor(n)
		root.swapNodes(n, pred)
	}

	doAssert(n.left == nil || n.right == nil)
	child := n.right
	if child == nil {
		child = n.left
	}
	if n.color == black {
		n.color = getColor(child)
		root.deleteCase1(n)
	}
	root.replaceNode(n, child)
	if n.parent == nil && child != nil {
		child.color = black
	}
	root.count--
	if root.count == 0 {
		root.minNode = nil
		root.maxNode = nil
	} else {
		if root.minNode == n {
			root.recomputeMinNode()
		}
		if root.maxNode == n {
			root.recomputeMaxNode()
		}
	}
}

// Move n to the pred's place, and vice versa
//
// TODO: this code is overly convoluted
func (root *rbTree[K, V]) swapNodes(n, pred *node[K, V]) {
	doAssert(pred != n)
	isLeft := pred.isLeftChild()
	tmp := *pred
	root.replaceNode(n, pred)
	pred.color = n.color

	if tmp.parent == n {
		// swap the positions of n and pred
		if isLeft {
			pred.left = n
			pred.right = n.right
			if pred.right != nil {
				pred.right.parent = pred
			}
		} else {
			pred.left = n.left
			if pred.left != nil {
				pred.left.parent = pred
			}
			pred.right = n
		}
		n.parent = pred
		n.left = tmp.left
		if n.left != nil {
			n.left.parent = n
		}
		n.right = tmp.right
		if n.right != nil {
			n.right.parent = n
		}
	} else {
		pred.left = n.left
		if pred.left != nil {
			pred.left.parent = pred
		}
		pred.right = n.right
		if pred.right != nil {
			pred.right.parent = pred
		}
		if isLeft {
			tmp.parent.left = n
		} else {
			tmp.parent.right = n
		}
		n.parent = tmp.parent
		n.left = tmp.left
		if n.left != nil {
			n.left.parent = n
		}
		n.right = tmp.right
		if n.right != nil {
			n.right.parent = n
		}
	}
	n.color = tmp.color
}

func (root *rbTree[K, V]) deleteCase1(n *node[K, V]) {
	for true {
		if n.parent != nil {
			if getColor(n.sibling()) == red {
				n.parent.color = red
				n.sibling().color = black
				if n == n.parent.left {
					root.rotateLeft(n.parent)
				} else {
					root.rotateRight(n.parent)
				}
			}
			if getColor(n.parent) == black &&
				getColor(n.sibling()) == black &&
				getColor(n.sibling().left) == black &&
				getColor(n.sibling().right) == black {
				n.sibling().color = red
				n = n.parent
				continue
			} else {
				// case 4
				if getColor(n.parent) == red &&
					getColor(n.sibling()) == black &&
					getColor(n.sibling().left) == black &&
					getColor(n.sibling().right) == black {
					n.sibling().color = red
					n.parent.color = black
				} else {
					root.deleteCase5(n)
				}
			}
		}
		break
	}
}

func (root *rbTree[K, V]) deleteCase5(n *node[K, V]) {
	if n == n.parent.left &&
		getColor(n.sibling()) == black &&
		getColor(n.sibling().left) == red &&
		getColor(n.sibling().right) == black {
		n.sibling().color = red
		n.sibling().left.color = black
		root.rotateRight(n.sibling())
	} else if n == n.parent.right &&
		getColor(n.sibling()) == black &&
		getColor(n.sibling().right) == red &&
		getColor(n.sibling().left) == black {
		n.sibling().color = red
		n.sibling().right.color = black
		root.rotateLeft(n.sibling())
	}

	// case 6
	n.sibling().color = getColor(n.parent)
	n.parent.color = black
	if n == n.parent.left {
		doAssert(getColor(n.sibling().right) == red)
		n.sibling().right.color = black
		root.rotateLeft(n.parent)
	} else {
		doAssert(getColor(n.sibling().left) == red)
		n.sibling().left.color = black
		root.rotateRight(n.parent)
	}
}

func (root *rbTree[K, V]) replaceNode(oldn, newn *node[K, V]) {
	if oldn.parent == nil {
		root.root = newn
	} else {
		if oldn.isLeftChild() {
			oldn.parent.left = newn
		} else {
			oldn.parent.right = newn
		}
	}
	if newn != nil {
		newn.parent = oldn.parent
	}
}

/*
    X		     Y
  A   Y	    => X   C
     B C 	  A B
*/
func (root *rbTree[K, V]) rotateLeft(n *node[K, V]) {
	r := n.right
	root.replaceNode(n, r)
	n.right = r.left
	if r.left != nil {
		r.left.parent = n
	}
	r.left = n
	n.parent = r

	/*
		y := x.right
		if y == nil {
			root.Dump()
			panic("about to crash b/c y is nil")
		}
		x.right = y.left
		if y.left != nil {
			y.left.parent = x
		}
		y.parent = x.parent
		if x.parent == nil {
			root.root = y
		} else {
			if x.isLeftChild() {
				x.parent.left = y
			} else {
				x.parent.right = y
			}
		}
		y.left = x
		x.parent = y
	*/
}

/*
     Y           X
   X   C  =>   A   Y
  A B             B C
*/
func (root *rbTree[K, V]) rotateRight(n *node[K, V]) {
	L := n.left
	root.replaceNode(n, L)
	n.left = L.right
	if L.right != nil {
		L.right.parent = n
	}
	L.right = n
	n.parent = L
}

func (root *rbTree[K, V]) DumpAsString() string {
	s := ""
	i := 0
	for nd := root.Min(); nd != nil; nd = nd.doNext() {
		s += fmt.Sprintf("node %03d: %#v\n", i, nd.item)
		i++
	}
	return s
}

func (root *rbTree[K, V]) Dump() {
	i := 0
	for nd := root.Min(); nd != nil; nd = nd.doNext() {
		fmt.Printf("node %03d: %#v\n", i, nd.item)
		i++
	}
	n := root.root
	for n.parent != nil {
		n = n.parent
	}

	root.Walk(n, 0, "root")
}

func colorString[K, V any](n *node[K, V]) string {
	if n.color == red {
		return "red"
	}
	return "black"
}

func (tr *rbTree[K, V]) Walk(n *node[K, V], indent int, lab string) {

	spc := strings.Repeat(" ", indent*3)
	var parItem, leftItem, rightItem interface{}
	if n.parent != nil {
		parItem = n.parent.item
	}
	if n.left != nil {
		leftItem = n.left.item
	}
	if n.right != nil {
		rightItem = n.right.item
	}
	fmt.Printf("%s %s node %p at indent %v [%s] %#v   leftChildNil:%v rightChildNil:%v.  my parent:'%#v'.  my left:'%#v', my right:'%#v'\n", spc, lab, n, indent, colorString(n), n.item, n.left == nil, n.right == nil, parItem, leftItem, rightItem)

	if n.left != nil {
		if n.left.parent != n {
			panic("n.left.parent != n")
		}
		tr.Walk(n.left, indent+1, "left")
	}

	if n.right != nil {
		if n.right.parent != n {
			panic("n.right.parent != n")
		}
		tr.Walk(n.right, indent+1, "right")
	}

	if n.color == red && n.parent.color == red {
		panic("double red chain found")
	}
}

var validations int

func validateTree2[K, V any](tr *rbTree[K, V]) {
	if tr == nil {
		panic("can't validate a nil tree")
	}
	root := tr.root
	if root == nil {
		return
	}
	for root.parent != nil {
		//vv("validateTree warning, not passed the root.")
		root = root.parent
	}
	tr.validateTreeHelper(root)
	//fmt.Printf("\n tree validated\n")
	validations++
}

func (tr *rbTree[K, V]) validateTreeHelper(n *node[K, V]) {

	if n.parent != nil {
		if n.parent.left != n && n.parent.right != n {
			panic("my parent doesn't know me")
		}
	}
	if n.left != nil {
		if n.left.parent != n {
			panic("my child doesn't know me")
		}
		tr.validateTreeHelper(n.left)
	}
	if n.right != nil {
		if n.right.parent != n {
			panic("my child doesn't know me")
		}
		tr.validateTreeHelper(n.right)
	}
}
