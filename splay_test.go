package gotree

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func init() {
}

var _ = spew.Dump

var _ = fmt.Printf

func (t *Tree) Height() int {
	var calc func(n *Node) int
	calc = func(n *Node) int {
		if n == nil {
			return 0
		}
		if a, b := calc(n.left), calc(n.right); a >= b {
			return 1 + a
		} else {
			return 1 + b
		}
	}
	return calc(t.root)
}

func TestSplayErrorInsert(t *testing.T) {
	t.Parallel()

	var old Interface
	var check error
	tree := &Tree{}

	old, check = tree.Insert(nil)
	if _, ok := check.(InvalidInterfaceError); !ok || old != nil {
		t.Errorf("Should Not be able to input nil")
		t.Errorf("Error should be of type InvalidInterfaceError")
	}
}

func TestSplayBasicInsert(t *testing.T) {

	t.Parallel()
	var old Interface
	var check error
	items := []exStruct{exStruct{0, "0"},
		exStruct{2, "2"}, exStruct{2, "3"}}
	tree := &Tree{}

	old, check = tree.Insert(items[0])
	if check != nil || old != nil {
		t.Errorf("First check on input is messed!")
	}

	tree.Insert(items[1])
	old, check = tree.Insert(items[2])
	if check != nil {
		t.Errorf("Check on old input is messed!")
	}
	if tree.Size != 2 {
		t.Errorf("Problems tracking Size")
	}
	if old != items[1] {
		t.Errorf("old input is messed!")
	}
	old, check = tree.Insert(items[1])
	if check != nil {
		t.Errorf("Check on old input is messed!")
	}
	if tree.Size != 2 {
		t.Errorf("Problems tracking Size")
	}
	if old != items[2] {
		t.Errorf("old input is messed!")
	}
}

func TestSplayInsert(t *testing.T) {

	var check error
	var old Interface
	tree := Tree{}
	for i := 0; i < iters; i++ {
		old, check = tree.Insert(exInt(i))
		if check != nil {
			t.Errorf("Check on old input is messed!")
		}
		if old != nil {
			t.Errorf("Old should be nil")
		}

	}
	//fmt.Println(tree.Height())
}

func TestSplaySearch(t *testing.T) {

	t.Parallel()
	tree := &Tree{}
	for i := 0; i < iters; i++ {
		tree.Insert(exInt(i))
	}
	_, ok := tree.Search(nil)
	if ok == nil {
		t.Errorf("Not minding nil key's")
	}

	//tree.Map(InOrder, inc(t))
	for i := 0; i < iters; i++ {
		value, ok := tree.Search(exInt(i))
		if ok != nil {
			t.Errorf("All these values should be present")
		}
		if int(value.(exInt)) != i {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}

	for i := iters; i < iters+1000; i++ {
		value, ok := tree.Search(exInt(i))
		if ok != nil {
			t.Errorf("values should not be present")
		}
		if value != nil {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}
}

func BenchmarkSplayInsert(b *testing.B) {

	b.StopTimer()
	tree := &Tree{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(exInt(b.N - i))
	}

}

func TestSplayIterIn(t *testing.T) {

	t.Parallel()
	tree := &Tree{}
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for _, v := range items {
		tree.Insert(exString(v))
	}
	var check *Node
	if check = tree.Next(); check != nil {
		t.Errorf("Didn't avoid a non intialized next call")
	}

	if tree.iterNext != nil {
		t.Errorf("Didn't reset iter")
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

	if check = tree.Next(); check != nil {
		t.Errorf("Didn't avoid a non intialized next call")
	}
	if tree.iterNext != nil {
		t.Errorf("Didn't reset iter")
	}

	//scs := spew.ConfigState{Indent: "\t"}
	//scs.Dump(tree.root)

}
