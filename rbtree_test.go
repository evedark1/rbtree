//
// Created by Yaz Saito on 06/10/12.
//

package rbtree

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"testing"
)

const testVerbose = false

// Create a tree storing a set of integers
func testNewIntTree() *rbTree[int, voidType] {
	return newTreeFunc[int, voidType](func(a, b int) int {
		return a - b
	})
}

func testAssert(t *testing.T, b bool, message string) {
	if !b {
		t.Fatal(message)
	}
}

func equalAssert[T comparable](t *testing.T, expect T, actual T) {
	if expect != actual {
		t.Fatal(fmt.Sprintf("equalAssert fail: %v != %v", expect, actual))
	}
}

func TestEmpty(t *testing.T) {
	tree := testNewIntTree()
	testAssert(t, tree.Len() == 0, "len!=0")
	testAssert(t, tree.Min() == nil, "minlimit")
	testAssert(t, tree.Max() == nil, "maxlimit")
	testAssert(t, tree.Find(10) == nil, "Not empty")
	testAssert(t, tree.FindGE(10) == nil, "Not empty")
	testAssert(t, tree.FindLE(10) == nil, "Not empty")
}

func TestFindGE(t *testing.T) {
	tree := testNewIntTree()
	tree.Insert(10)
	testAssert(t, tree.Len() == 1, "len==1")
	testAssert(t, tree.FindGE(10).item == 10, "FindGE 10")
	testAssert(t, tree.FindGE(11) == nil, "FindGE 11")
	testAssert(t, tree.FindGE(9).item == 10, "FindGE 10")
}

func TestFindLE(t *testing.T) {
	tree := testNewIntTree()
	tree.Insert(10)
	testAssert(t, tree.FindLE(10).item == 10, "FindLE 10")
	testAssert(t, tree.FindLE(11).item == 10, "FindLE 11")
	testAssert(t, tree.FindLE(9) == nil, "FindLE 9")
}

func TestTreeFind(t *testing.T) {
	tree := testNewIntTree()
	_, ok := tree.Insert(10)
	testAssert(t, ok, "insert1")
	_, ok = tree.Insert(10)
	testAssert(t, !ok, "insert1")
	testAssert(t, tree.Find(10).item == 10, "Get 10")
	testAssert(t, tree.Find(9) == nil, "Get 9")
	testAssert(t, tree.Find(11) == nil, "Get 11")
}

func TestDelete(t *testing.T) {
	tree := testNewIntTree()
	testAssert(t, !tree.DeleteWithKey(10), "del")
	testAssert(t, tree.Len() == 0, "dellen")
	tree.Insert(10)
	testAssert(t, tree.DeleteWithKey(10), "del")
	testAssert(t, tree.Len() == 0, "dellen")

	// delete was deleting after the request if request not found
	// ensure this does not regress:
	tree.Insert(10)
	testAssert(t, !tree.DeleteWithKey(9), "del")
	testAssert(t, tree.Len() == 1, "dellen")

}

func iterToString(nd *node[int, voidType]) string {
	s := ""
	for ; nd != nil; nd = nd.doNext() {
		if s != "" {
			s = s + ","
		}
		s = s + fmt.Sprintf("%d", nd.item)
	}
	return s
}

func reverseIterToString(nd *node[int, voidType]) string {
	s := ""
	for ; nd != nil; nd = nd.doPrev() {
		if s != "" {
			s = s + ","
		}
		s = s + fmt.Sprintf("%d", nd.item)
	}
	return s
}

func TestIterator(t *testing.T) {
	tree := testNewIntTree()
	for i := 0; i < 10; i = i + 2 {
		tree.Insert(i)
	}
	if iterToString(tree.FindGE(3)) != "4,6,8" {
		t.Error("iter")
	}
	if iterToString(tree.FindGE(4)) != "4,6,8" {
		t.Error("iter")
	}
	if iterToString(tree.FindGE(8)) != "8" {
		t.Error("iter")
	}
	if iterToString(tree.FindGE(9)) != "" {
		t.Error("iter")
	}

	if reverseIterToString(tree.FindLE(3)) != "2,0" {
		t.Error("iter", reverseIterToString(tree.FindLE(3)))
	}
	if reverseIterToString(tree.FindLE(2)) != "2,0" {
		t.Error("iter")
	}
	if reverseIterToString(tree.FindLE(0)) != "0" {
		t.Error("iter")
	}
	if reverseIterToString(tree.FindLE(-1)) != "" {
		t.Error("iter")
	}
}

//
// Randomized tests
//

// oracle stores provides an interface similar to rbtree, but stores
// data in an sorted array
type oracle struct {
	data []int
}

func newOracle() *oracle {
	return &oracle{data: make([]int, 0)}
}

func (o *oracle) Len() int {
	return len(o.data)
}

// interface needed for sorting
func (o *oracle) Less(i, j int) bool {
	return o.data[i] < o.data[j]
}

func (o *oracle) Swap(i, j int) {
	e := o.data[j]
	o.data[j] = o.data[i]
	o.data[i] = e
}

func (o *oracle) Insert(key int) bool {
	for _, e := range o.data {
		if e == key {
			return false
		}
	}

	o.data = append(o.data, key)
	sort.Sort(o)
	return true
}

func (o *oracle) RandomExistingKey(rand *rand.Rand) int {
	index := rand.Intn(len(o.data))
	return o.data[index]
}

func (o *oracle) FindGE(t *testing.T, key int) oracleIterator {
	prev := int(-1)
	for i, e := range o.data {
		if e <= prev {
			t.Fatal("Nonsorted oracle ", e, prev)
		}
		if e >= key {
			return oracleIterator{o: o, index: i}
		}
	}
	return oracleIterator{o: o, index: len(o.data)}
}

func (o *oracle) FindLE(t *testing.T, key int) oracleIterator {
	iter := o.FindGE(t, key)
	if !iter.Limit() && o.data[iter.index] == key {
		return iter
	}
	return oracleIterator{o, iter.index - 1}
}

func (o *oracle) Delete(key int) bool {
	for i, e := range o.data {
		if e == key {
			newData := make([]int, len(o.data)-1)
			copy(newData, o.data[0:i])
			copy(newData[i:], o.data[i+1:])
			o.data = newData
			return true
		}
	}
	return false
}

//
// Test iterator
//
type oracleIterator struct {
	o     *oracle
	index int
}

func (oiter oracleIterator) Limit() bool {
	return oiter.index >= len(oiter.o.data) || oiter.index < 0
}

func (oiter oracleIterator) Min() bool {
	return oiter.index == 0
}

func (oiter oracleIterator) Max() bool {
	return oiter.index == len(oiter.o.data)-1
}

func (oiter oracleIterator) Item() int {
	return oiter.o.data[oiter.index]
}

func (oiter oracleIterator) Next() oracleIterator {
	return oracleIterator{oiter.o, oiter.index + 1}
}

func (oiter oracleIterator) Prev() oracleIterator {
	return oracleIterator{oiter.o, oiter.index - 1}
}

func compareContents(t *testing.T, oiter oracleIterator, titer *node[int, voidType]) {
	oi := oiter
	ti := titer

	// Test forward iteration
	testAssert(t, oi.Limit() == (ti == nil), "rend")
	if oi.Limit() {
		return
	}

	for !oi.Limit() && ti != nil {
		// log.Print("Item: ", oi.Item(), ti.Item())
		if ti.item != oi.Item() {
			t.Fatal("Wrong item", ti.item, oi.Item())
		}
		oi = oi.Next()
		ti = ti.doNext()
	}
	if ti != nil {
		t.Fatal("!ti.done", ti.item)
	}
	if !oi.Limit() {
		t.Fatal("!oi.done", oi.Item())
	}

	// Test reverse iteration
	oi = oiter
	ti = titer

	for !oi.Limit() && ti != nil {
		if ti.item != oi.Item() {
			t.Fatal("Wrong item", ti.item, oi.Item())
		}
		oi = oi.Prev()
		ti = ti.doPrev()
	}
	if ti != nil {
		t.Fatal("!ti.done", ti.item)
	}
	if !oi.Limit() {
		t.Fatal("!oi.done", oi.Item())
	}
}

func compareContentsFull(t *testing.T, o *oracle, tree *rbTree[int, voidType]) {
	compareContents(t, o.FindGE(t, int(-1)), tree.FindGE(-1))
}

func TestRandomized(t *testing.T) {
	const numKeys = 1000

	o := newOracle()
	tree := testNewIntTree()
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 10000; i++ {
		op := r.Intn(100)
		if op < 50 {
			key := r.Intn(numKeys)
			if testVerbose {
				log.Print("Insert ", key)
			}
			o.Insert(key)
			tree.Insert(key)
			compareContentsFull(t, o, tree)
		} else if op < 90 && o.Len() > 0 {
			key := o.RandomExistingKey(r)
			if testVerbose {
				log.Print("DeleteExisting ", key)
			}
			o.Delete(key)
			if !tree.DeleteWithKey(key) {
				t.Fatal("DeleteExisting", key)
			}
			compareContentsFull(t, o, tree)
		} else if op < 95 {
			key := int(r.Intn(numKeys))
			if testVerbose {
				log.Print("FindGE ", key)
			}
			compareContents(t, o.FindGE(t, key), tree.FindGE(key))
		} else {
			key := int(r.Intn(numKeys))
			if testVerbose {
				log.Print("FindLE ", key)
			}
			compareContents(t, o.FindLE(t, key), tree.FindLE(key))
		}
	}
}
