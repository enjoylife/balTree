package gotree

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math/rand"
	"testing"
)

var _ = spew.Dump

var _ = fmt.Printf

/* Testing Compare function: int */
func testCmpInt(a interface{}, b interface{}) Direction {
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
func testCmpString(c interface{}, d interface{}) Direction {
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
		if items[i] != n.Value() {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value(), items[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	for i, n := 0, tree.InitIter(PreOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if preOrder[i] != n.Value() {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value(), preOrder[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	for i, n := 0, tree.InitIter(PostOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if postOrder[i] != n.Value() {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Value(), postOrder[i])
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

/* Experiments: not used */
func (t *RbTree) insertIter(h *rbNode, key interface{}, value interface{}) *rbNode {

	// empty tree
	if h == nil {
		return &rbNode{color: Red, key: key, value: value}
	}

	// setup our own stack and helpers
	var (
		stack     = []*rbNode{}
		count int = 0
		prior *rbNode
	)

L:
	for {
		switch t.cmp(h.Key, key) {
		case EQ:
			h.value = value
			return t.root // no need for rest of the fix code
		case LT:
			prior = h
			stack = append(stack, prior)
			count++
			h = h.left
			if h == nil {
				h = &rbNode{color: Red, key: key, value: value}
				prior.left = h
				break L
			}
		case GT:
			prior = h
			stack = append(stack, prior)
			count++
			h = h.right
			if h == nil {
				h = &rbNode{color: Red, key: key, value: value}
				prior.right = h
				break L
			}
		default:
			panic("Compare result undefined")
		}

		if prior == h {
			panic("Shouldn't be equal last check")
		}

	}

	// h is parent of new node at this point
	h = prior
L2:
	for {
		count--

		if h.right.isRed() && !(h.left.isRed()) {
			h = h.rotateLeft()
		}
		if h.left.isRed() && h.left.left.isRed() {
			h = h.rotateRight()
		}
		if h.left.isRed() && h.right.isRed() {
			h.colorFlip()
		}

		if count == 0 {
			break L2
		}

		if count > 0 {

			switch t.cmp(stack[count-1].Key, h.Key) {
			case LT:
				stack[count-1].left = h
			case GT:
				stack[count-1].right = h
			}
			h = stack[count-1]
		}

	}

	return h
}
func (t *RbTree) TraverseIter(f IterFunc) {
	node := t.root
	stack := []*rbNode{nil}
	for len(stack) != 1 || node != nil {
		if node != nil {
			stack = append(stack, node)
			node = node.left
		} else {
			stackIndex := len(stack) - 1
			node = stack[stackIndex]
			f(node.Key, node.value)
			stack = stack[0:stackIndex]
			node = node.right
		}
	}
}
