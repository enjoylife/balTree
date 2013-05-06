package gotree

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math/rand"
	"testing"
)

func init() {
	// so that every run is the same seq of rand numbers
	rand.Seed(0)
}

var _ = spew.Dump

var _ = fmt.Printf

const (
	searchTotal = 100000
	searchSpace = searchTotal / 2
)

type ExInt int

func (this ExInt) Compare(b Comparable) Balance {
	switch that := b.(type) {
	case ExInt:
		switch result := int(this - that); {
		case result > 0:
			return LT
		case result < 0:
			return GT
		case result == 0:
			return EQ
		default:
			panic("Invalid Compare function Result")
		}

	default:
		s := fmt.Sprintf("Can not compare to the unkown type of %T", that)
		panic(s)
	}

}

/* Testing Compare function: int */
func testCmpInt(a interface{}, b interface{}) Balance {
	switch result := (a.(int) - b.(int)); {
	case result > 0:
		return LT
	case result < 0:
		return GT
	case result == 0:
		return EQ
	default:
		panic("Invalid Compare function Result")
	}
}

/* Testing Compare function: string */
func testCmpString(c interface{}, d interface{}) Balance {
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
		return LT
	case result < 0:
		return GT
	case result == 0:
		return EQ
	default:
		panic("Invalid Compare function Result")
	}
}

/* Helpers for tree traversal and testing tree properties */

func printNode(n *Node) {
	fmt.Println("VALUE: ", n.Value)
}

func isBalanced(t *RBTree) bool {
	if t == nil {
		return true
	}
	var black int // number of black links on path from root to min
	black = 0
	for x := t.root; x != nil; x = x.left {
		if x.color == Black {
			black++
		}
	}
	return nodeIsBalanced(t.root, black) && t.Height == black
}

func nodeIsBalanced(n *Node, black int) bool {
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

func inc(t *testing.T) func(n *Node) {
	var prior int = -1
	return func(n *Node) {
		if prior < n.Value.(int) {
			//fmt.Println("VALUE: ", value.(int))
			prior = n.Value.(int)
		} else {
			t.Errorf("Prior: %d, Current: %d", prior, n.Value.(int))
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
	tree.Traverse(InOrder, inc(t))
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

	tree.Traverse(InOrder, inc(t))
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

func TestRemove(t *testing.T) {

	tree := New(testCmpInt)
	iters := 1000
	for i := 0; i < iters; i++ {
		tree.Insert(i, i)
	}

	for i := 0; i < iters; i++ {
		tree.Remove(i)

		black := 0
		for x := tree.root; x != nil; x = x.left {
			if x.color == Black {
				black++
			}
		}

		if !isBalanced(tree) {
			fmt.Println("Height", tree.Height)
			fmt.Println("Calc Height", black)
			t.Errorf("Tree is not balanced")
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

	for i, n := 0, tree.InitIter(InOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if items[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, items[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	for i, n := 0, tree.InitIter(PreOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if preOrder[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, preOrder[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	for i, n := 0, tree.InitIter(PostOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if postOrder[i] != n.Value {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value, postOrder[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}

	//tree.Traverse(PreOrder, printNode)
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
	//tree.Traverse(InOrder, printNode)

}

func BenchmarkMapInsert(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	m := make(map[int]int)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := r.Int()
		m[a] = a
	}

}

func BenchmarkInsert(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := r.Int()
		tree.Insert(a, a)
	}
}
func BenchmarkSearch(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	m := make(map[int]int)
	tree := New(testCmpInt)
	for i := 0; i < searchTotal; i++ {
		a := r.Intn(searchSpace)
		m[a] = a
		tree.Insert(a, a)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := r.Intn(searchSpace)
		tree.Search(a)
	}
}

func BenchmarkRemove(b *testing.B) {

	b.StopTimer()
	tree := New(testCmpInt)
	for i := 0; i < b.N; i++ {
		tree.Insert((b.N - i), i)
	}
	b.StartTimer()
	/*for i := 0; i < b.N; i++ {
		a := r.Intn(searchSpace)
		tree.Remove(a)
	}*/
	for i := 0; i < b.N; i++ {
		tree.Remove(i)
	}
}

func BenchmarkIterInOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		for i, n := 0, tree.InitIter(InOrder); n != nil; i, n = i+1, tree.Next() {
			sum += n.Value.(int)
		}
	}

}
func BenchmarkIterPreOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		for i, n := 0, tree.InitIter(PreOrder); n != nil; i, n = i+1, tree.Next() {
			sum += n.Value.(int)
		}
	}

}
func BenchmarkIterPostOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		for i, n := 0, tree.InitIter(PostOrder); n != nil; i, n = i+1, tree.Next() {
			sum += n.Value.(int)
		}
	}

}

func recurse() func(n *Node) {
	var sum int = 0
	return func(n *Node) {
		sum += n.Value.(int)
	}
}
func BenchmarkRecurseTraverseInorderOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f := recurse()
		tree.Traverse(InOrder, f)

	}

}
func BenchmarkRecurseTraversePreorderOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f := recurse()
		tree.Traverse(PreOrder, f)

	}

}
func BenchmarkRecurseTraversePostOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmpInt)
	for i := 0; i < 1000; i++ {
		tree.Insert(r.Int(), r.Int())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f := recurse()
		tree.Traverse(PostOrder, f)

	}

}
