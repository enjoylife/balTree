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

func TestSplayBasicInsert(t *testing.T) {

	var old Interface
	items := []exStruct{exStruct{0, "0"},
		exStruct{2, "2"}, exStruct{2, "3"}}
	tree := &SplayTree{}

	old = tree.Insert(items[0])
	if old != nil {
		t.Errorf("First check on input is messed!")
	}

	tree.Insert(items[1])
	old = tree.Insert(items[2])
	if tree.Size != 2 {
		t.Errorf("Problems tracking Size")
	}
	if old != items[1] {
		t.Errorf("old input is messed!")
	}
	old = tree.Insert(items[1])
	if tree.Size != 2 {
		t.Errorf("Problems tracking Size")
	}
	if old != items[2] {
		t.Errorf("old input is messed!")
	}
}

func TestSplayInsert(t *testing.T) {

	var old Interface
	tree := SplayTree{}
	for i := 0; i < iters; i++ {
		old = tree.Insert(exInt(i))
		if old != nil {
			t.Errorf("Old should be nil")
		}
		old = tree.Insert(exInt(i))
		if old != exInt(i) {
			t.Errorf("Old should match")
		}

	}
	//fmt.Println(tree.Height())
}

func TestSplaySearch(t *testing.T) {

	tree := &SplayTree{}

	x := tree.Search(nil)
	if x != nil {
		t.Errorf("Not minding empty tree")
	}
	for i := 0; i < iters; i++ {
		tree.Insert(exInt(i))
	}
	x = tree.Search(nil)
	if x != nil {
		t.Errorf("Not minding nil key's")
	}

	for i := 0; i < iters; i++ {
		value := tree.Search(exInt(i))
		if int(value.(exInt)) != i {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}

	for i := iters; i < iters+1000; i++ {
		value := tree.Search(exInt(i))
		if value != nil {
			t.Errorf("Values don't match Exp: %d, Got: %d", i, value)
		}
	}
}

func BenchmarkSplayInsert(b *testing.B) {

	b.StopTimer()
	tree := &SplayTree{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(exInt(b.N - i))
	}

}

func TestSplayRemove(t *testing.T) {

	tree := &SplayTree{}
	c := tree.Remove(exInt(0))
	if c != nil {
		t.Errorf("Not respecting empty tree.")
	}
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for _, v := range items {
		tree.Insert(exString(v))
	}
	for k, v := range items {
		check := tree.Remove(exString(v))
		if check != exString(v) {
			fmt.Println(check)
			t.Errorf("Not getting removed item back")
		}
		check = tree.Search(exString(v))
		if check != nil {
			t.Errorf("Didn't really remove")
		}

		for i, n := k+1, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {
			if items[i] != string(n.Elem.(exString)) {
				t.Errorf("Other elems deleted Got:%s, Exp: %s", n.Elem, items[i])
			}
		}
	}
}
func TestSplayMaxRemove(t *testing.T) {
	t.Parallel()
	tree := &RBTree{}

	for i := 0; i < iters; i++ {
		tree.Insert(exInt(i))
		if tree.Max() != exInt(i) {
			t.Errorf("Max not updateing")
		}
	}
	for i := iters; i > 0; i-- {
		tree.Remove(exInt(i))
		if tree.Max() != exInt(i-1) {
			fmt.Println(tree.Max())
			t.Errorf("Max not updateing")
		}
	}
	tree.Remove(exInt(0))
	if tree.Max() != nil {
		fmt.Println(tree.Max())
		t.Errorf("Max not updateing")
	}

}
func TestSplayMinRemove(t *testing.T) {
	t.Parallel()
	tree := &RBTree{}
	var old Interface
	var check error

	for i := iters; i >= 0; i-- {
		tree.Insert(exInt(i))
		if tree.Min() != exInt(i) {
			fmt.Println(tree.Min())
			t.Errorf("Min not updateing")
		}
	}
	for i := 0; i < iters; i++ {

		old = tree.Remove(exInt(i))
		if old == nil {
			fmt.Println("old", old)
			fmt.Println(check)
			t.Errorf("Not giving back old value")
		}
		if tree.Min() != exInt(i+1) {
			fmt.Println(tree.Min())
			t.Errorf("Min not updateing")
		}
	}

	tree.Remove(exInt(iters))
	if tree.Min() != nil {
		fmt.Println(tree.Min())
		t.Errorf("Min not working")
	}
}
func TestSplaySizeRemove(t *testing.T) {
	t.Parallel()
	tree := &RBTree{}

	for i := iters; i >= 0; i-- {
		tree.Insert(exInt(i))
	}

	for i := 0; i < iters; i++ {
		tree.Remove(exInt(i))
		if tree.Size != (iters - i) {
			fmt.Println(tree.Size)
			t.Errorf("Size on remove not working")
		}

	}
}

func TestSplayIterIn(t *testing.T) {

	tree := &SplayTree{}
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for _, v := range items {
		tree.Insert(exString(v))
	}
	var check *SplayNode
	if check = tree.Next(); check != nil {
		t.Errorf("Didn't avoid a non intialized next call")
	}

	if tree.iterNext != nil {
		t.Errorf("Didn't reset iter")
	}
	count := 0

	for i, n := 0, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {

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

}
