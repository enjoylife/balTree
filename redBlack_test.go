package gotree

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math/rand"
	"testing"
)

func init() {
}

var _ = spew.Dump

var _ = fmt.Printf

const (
	searchTotal = 100000
	searchSpace = searchTotal / 2
	iters       = 10000
)

type exInt int

func (this exInt) Compare(b Comparer) Balance {
	switch that := b.(type) {
	case exInt:
		switch result := int(this - that); {
		case result > 0:
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}
	default:
		return NP
		s := fmt.Sprintf("Can not compare to the unkown type of %T", that)
		panic(s)
	}

}

type exString string

func (this exString) Compare(b Comparer) Balance {
	switch that := b.(type) {
	case exString:
		a := string(this)
		b := string(that)
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
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}

	default:
		return NP
		s := fmt.Sprintf("Can not compare to the unkown type of %T", that)
		panic(s)
	}
}

type exStruct struct {
	M int
	S string
}

func (this exStruct) Compare(b Comparer) Balance {
	switch that := b.(type) {
	case exStruct:
		switch result := int(this.M - that.M); {
		case result > 0:
			return GT
		case result < 0:
			return LT
		case result == 0:
			return EQ
		default:
			return NP
		}
	default:
		return NP
		s := fmt.Sprintf("Can not compare to the unkown type of %T", that)
		panic(s)
	}
}

/* Helpers for tree traversal and testing tree properties */
func printNode(n *Node) {
	x := n.Elem
	fmt.Println("Elem:", x)
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

// can only be used for exInt
func inc(t *testing.T) func(n *Node) {
	var prior int = -1
	return func(n *Node) {
		if prior < int(n.Elem.(exInt)) {
			//fmt.Println("VALUE: ", value.(int))
			prior = int(n.Elem.(exInt))
		} else {
			t.Errorf("Prior: %d, Current: %d", prior, n.Elem)
		}
	}
}
func TestBasicInsert(t *testing.T) {

	var old Comparer
	var check error
	items := []exStruct{exStruct{0, "0"},
		exStruct{2, "2"}, exStruct{2, "3"}}
	tree := &RBTree{}
	old, check = tree.Insert(nil)
	if old != nil || check == nil {
		t.Errorf("Should Not be able to input nil")
	}

	if tree.Min() != nil {
		fmt.Println(tree.Min())
		t.Errorf("Min not working")
	}
	if tree.Max() != nil {
		fmt.Println(tree.Max())
		t.Errorf("Max not working")
	}
	old, check = tree.Insert(items[0])
	if check != nil || old != nil {
		t.Errorf("First check on input is messed!")
	}

	tree.Insert(items[1])
	old, check = tree.Insert(items[2])
	if check != nil {
		t.Errorf("Check on old input is messed!")
	}
	if old != items[1] {
		t.Errorf("old input is messed!")
	}
	old, check = tree.Insert(items[1])
	if check != nil {
		t.Errorf("Check on old input is messed!")
	}
	if old != items[2] {
		t.Errorf("old input is messed!")
	}
	if tree.Min() != items[0] {
		fmt.Println(tree.Min())
		t.Errorf("Min not working")
	}
	if tree.Max() != items[1] {
		fmt.Println(tree.Max())
		t.Errorf("Max not working")
	}
}

func TestMoreInsert(t *testing.T) {

	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}
	for i := 0; i < iters; i++ {
		item := r.Int()
		tree.Insert(exInt(item))
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}
	tree.Traverse(InOrder, inc(t))
}

func TestSearch(t *testing.T) {

	tree := &RBTree{}
	for i := 0; i < iters; i++ {
		tree.Insert(exInt(i))
	}
	_, ok := tree.Search(nil)
	if ok == nil {
		t.Errorf("Not minding nil key's")
	}

	tree.Traverse(InOrder, inc(t))
	for i := 0; i < iters; i++ {
		value, ok := tree.Search(exInt(i))
		if ok != nil {
			t.Errorf("All these values should be present")
		}
		if int(value.(exInt)) != i {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}

	for i := iters; i < iters*2; i++ {
		value, ok := tree.Search(exInt(i))
		if ok != nil {
			t.Errorf("values should not be present")
		}
		if value != nil {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}
}

func TestRemove(t *testing.T) {

	tree := &RBTree{}
	for i := 0; i < iters; i++ {
		tree.Insert(exInt(i))
	}

	for i := 0; i < iters; i++ {
		tree.Remove(exInt(i))

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

	tree := &RBTree{}
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	preOrder := []string{"d", "b", "a", "c", "h", "f", "e", "g", "i"}
	postOrder := []string{"a", "c", "b", "e", "g", "f", "i", "h", "d"}
	for _, v := range items {
		tree.Insert(exString(v))
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}

	count := 0

	for i, n := 0, tree.InitIter(InOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if items[i] != string(n.Elem.(exString)) {
			t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n.Elem, items[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	for i, n := 0, tree.InitIter(PreOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if preOrder[i] != string(n.Elem.(exString)) {
			t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n.Elem, preOrder[i])
		}

	}
	if count != len(items) {
		t.Errorf("Did not traverse all elements missing: %d", len(items)-count)
	}
	count = 0
	for i, n := 0, tree.InitIter(PostOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if postOrder[i] != string(n.Elem.(exString)) {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Elem, postOrder[i])
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

	tree := &RBTree{}
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for _, v := range items {
		tree.Insert(exString(v))
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}
	tree.Traverse(InOrder, printNode)

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
	tree := &RBTree{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(exInt(i))
	}
}
func BenchmarkSearch(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	m := make(map[int]int)
	tree := &RBTree{}
	for i := 0; i < searchTotal; i++ {
		a := r.Intn(searchSpace)
		m[a] = a
		tree.Insert(exInt(a))
		//tree.Insert(a, a)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := r.Intn(searchSpace)
		tree.Search(exInt(a))
	}
}

func BenchmarkRemove(b *testing.B) {

	b.StopTimer()
	tree := &RBTree{}
	for i := 0; i < b.N; i++ {
		tree.Insert(exInt(b.N - i))
	}
	b.StartTimer()
	/*for i := 0; i < b.N; i++ {
		a := r.Intn(searchSpace)
		tree.Remove(a)
	}*/
	for i := 0; i < b.N; i++ {
		tree.Remove(exInt(i))
	}
}

func BenchmarkIterInOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		for i, n := 0, tree.InitIter(InOrder); n != nil; i, n = i+1, tree.Next() {
			sum += int(n.Elem.(exInt))
		}
	}

}
func BenchmarkIterPreOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		for i, n := 0, tree.InitIter(PreOrder); n != nil; i, n = i+1, tree.Next() {
			sum += int(n.Elem.(exInt))
		}
	}

}
func BenchmarkIterPostOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum := 0

		for i, n := 0, tree.InitIter(PostOrder); n != nil; i, n = i+1, tree.Next() {
			sum += int(n.Elem.(exInt))
		}
	}

}

func recurse() func(n *Node) {
	var sum int = 0
	return func(n *Node) {
		sum += int(n.Elem.(exInt))
	}
}
func BenchmarkRecurseTraverseInorderOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
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
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
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
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f := recurse()
		tree.Traverse(PostOrder, f)

	}

}
