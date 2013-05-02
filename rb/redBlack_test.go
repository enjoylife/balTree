package gorbtree

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gotree"
	"math/rand"
	"testing"
)

var _ = spew.Dump

var _ = fmt.Printf

/* Testing Compare function: int */
func testCmpInt(a interface{}, b interface{}) gotree.Direction {
	switch result := (a.(int) - b.(int)); {
	case result > 0:
		return gotree.LT
	case result < 0:
		return gotree.GT
	case result == 0:
		return gotree.EQ
	default:
		panic("Invalid Compare function Result")
	}
}

/* Testing Compare function: string */
func testCmpString(c interface{}, d interface{}) gotree.Direction {
	a := c.(string)
	b := d.(string)
	min := len(b)
	if len(a) < len(b) {
		min = len(a)
	}
	diff := 0
	for i := 0; i < min && diff == 0; i++ {
		diff = int(a[i]) - int(b[i])
	}
	if diff == 0 {
		diff = len(a) - len(b)
	}

	switch result := diff; {
	case result > 0:
		return gotree.LT
	case result < 0:
		return gotree.GT
	case result == 0:
		return gotree.EQ
	default:
		panic("Invalid Compare function Result")
	}
}

/* Helpers for tree traversal and testing tree properties */

func printNode(key interface{}, value interface{}) {
	fmt.Println("VALUE: ", value)
}

func isBalanced(t *RbTree) bool {
	if t == nil {
		return true
	}
	var black int // number of black links on path from root to min
	for x := t.root; x != nil; x = x.left {
		if x.color == Black {
			black++
		}
	}
	return nodeIsBalanced(t.root, black) && t.Height == black
}

func nodeIsBalanced(n *rbNode, black int) bool {
	if n == nil && black == 0 {
		return true
	} else if n == nil && black != 0 {
		return false
	}
	if n.color == Black {
		black--
	}
	return nodeIsBalanced(n.left, black) && nodeIsBalanced(n.right, black)
}

func inc(t *testing.T) func(key interface{}, value interface{}) {
	var prior int = -1
	return func(key interface{}, value interface{}) {
		if prior < value.(int) {
			//fmt.Println("VALUE: ", value.(int))
			prior = value.(int)
		} else {
			t.Errorf("Prior: %d, Current: %d", prior, value.(int))
		}
	}
}

func TestInsert(t *testing.T) {

	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	iters := 10000
	for i := 0; i < iters; i++ {
		item := r.Int()
		tree.Insert(item, item)
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}
	tree.Traverse(gotree.InOrder, inc(t))
}

func TestSearch(t *testing.T) {

	tree := New(testCmpInt)
	iters := 10000
	for i := 0; i < iters; i++ {
		tree.Insert(i, i)
	}
	_, ok := tree.Search(nil)
	if ok {
		t.Errorf("Not minding nil key's")
	}

	tree.Traverse(gotree.InOrder, inc(t))
	for i := 0; i < iters; i++ {
		value, ok := tree.Search(i)
		if !ok {
			t.Errorf("All these values should be present")
		}
		if value != i {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}

	for i := iters; i < iters*2; i++ {
		value, ok := tree.Search(i)
		if ok {
			t.Errorf("values should not be present")
		}
		if value != nil {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}
}

func TestIterIn(t *testing.T) {

	tree := New(testCmpInt)
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	preOrder := []string{"d", "b", "a", "c", "h", "f", "e", "g", "i"}
	postOrder := []string{"a", "c", "b", "e", "g", "f", "i", "h", "d"}
	for i, v := range items {
		tree.Insert(i, v)
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}

	count := 0

	for i, n := 0, tree.InitIter(gotree.InOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if items[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, items[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	for i, n := 0, tree.InitIter(gotree.PreOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if preOrder[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, preOrder[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	for i, n := 0, tree.InitIter(gotree.PostOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if postOrder[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, postOrder[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}

	//tree.Traverse(gotree.PreOrder, printNode)
	//scs := spew.ConfigState{Indent: "\t"}
	//scs.Dump(tree.root)

}

func TestTraversal(t *testing.T) {
	tree := New(testCmpString)
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for _, v := range items {
		tree.Insert(v, v)
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}
	//tree.Traverse(gotree.InOrder, printNode)

}

func BenchmarkMapInsert(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	m := make(map[int]int)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[r.Int()] = r.Int()
	}

}

func BenchmarkInsert(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(r.Int(), r.Int())
	}

}
