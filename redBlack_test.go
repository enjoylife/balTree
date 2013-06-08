package gotree

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func isBalanced(t *RBTree) bool {
	if t == nil {
		return true
	}
	h := 0
	for x := t.root; x != nil; x = x.left {
		if x.color == black {
			h++
		}
	}
	return nodeIsBalanced(t.root, h) && t.Height() == h
}

func nodeIsBalanced(n *RBNode, h int) bool {
	if n == nil && h == 0 {
		return true
	} else if n == nil && h != 0 {
		return false
	}
	if n.color == black {
		h--
	}
	return nodeIsBalanced(n.left, h) && nodeIsBalanced(n.right, h)
}

// remove, min and max tests
func TestRBRemove(t *testing.T) {

	var old Interface
	tree := &RBTree{}

	old = tree.Remove(nil)
	if old != nil {
		t.Errorf("Should Not be able to remove nil")
	}

	item1 := exStruct{0, "0"}
	old = tree.Remove(exStruct{0, "1"})
	if old != nil {
		fmt.Println(old)
		fmt.Println(item1)
		t.Errorf("Not minding empty tree.")
	}

	tree.Insert(item1)
	old = tree.Remove(exStruct{0, "1"})
	if old != item1 {
		fmt.Println(old)
		fmt.Println(item1)
		t.Errorf("Can't even remove simple root")
	}
	old = tree.Search(exStruct{0, "1"})
	if old != nil {
		t.Errorf("Did not actually remove")
	}

	max := 100
	for i := 0; i < max; i++ {
		tree.Insert(exStruct{i, strconv.Itoa(i)})
	}

	for i := 0; i < max; i++ {
		tree.Remove(exStruct{i, strconv.Itoa(i)})
		old = tree.Search(exStruct{0, "1"})
		if old != nil {
			t.Errorf("Did not actually remove")
		}
	}
	for i := max; i < max*2; i++ {
		old = tree.Remove(exStruct{i, strconv.Itoa(i)})
		if old != nil {
			fmt.Println(old)
			t.Errorf("Can't  ignore nonexisitant elements in remove.")
		}
		h := 0
		for x := tree.root; x != nil; x = x.left {
			if x.color == black {
				h++
			}
		}

		if !isBalanced(tree) {
			fmt.Println("Height", tree.Height())
			fmt.Println("Calc Height", black)
			t.Errorf("Tree is not balanced")
		}
	}
}

func TestRBRandomRemove(t *testing.T) {

	tree := &RBTree{}
	r := rand.New(rand.NewSource(int64(5)))
	m := make(map[int]int)
	for i := 0; i < iters; i++ {
		a := r.Intn(searchSpace)
		m[a] = a
		tree.Insert(exInt(a))
	}

	for _, value := range m {
		tree.Remove(exInt(value))
		h := 0
		for x := tree.root; x != nil; x = x.left {
			if x.color == black {
				h++
			}
		}

		if !isBalanced(tree) {
			fmt.Println("Height", tree.Height())
			fmt.Println("Calc Height", h)
			t.Errorf("Tree is not balanced")
		}
	}
}

// iteration and map tests
func TestIterRBIn(t *testing.T) {

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
	var check Interface
	if check = tree.Next(); check != nil {
		t.Errorf("Didn't avoid a non intialized next call")
	}

	if tree.iterNext != nil {
		t.Errorf("Didn't reset iter")
	}
	count := 0

	for i, n := 0, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if items[i] != string(n.(exString)) {
			t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n, items[i])
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
	count = 0
	for i, n := 0, tree.IterInit(PreOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if preOrder[i] != string(n.(exString)) {
			t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n, preOrder[i])
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
	count = 0
	for i, n := 0, tree.IterInit(PostOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if postOrder[i] != string(n.(exString)) {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n, postOrder[i])
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

		for i, n := 0, tree.IterInit(PreOrder); n != nil; i, n = i+1, tree.Next() {
			sum += int(n.(exInt))
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

		for i, n := 0, tree.IterInit(PostOrder); n != nil; i, n = i+1, tree.Next() {
			sum += int(n.(exInt))
		}
	}

}

func BenchmarkMapInOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f := recurse()
		tree.Map(InOrder, f)

	}

}
func BenchmarkMapPreOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f := recurse()
		tree.Map(PreOrder, f)

	}

}
func BenchmarkMapPostOrder(b *testing.B) {

	b.StopTimer()
	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}
	for i := 0; i < 1000; i++ {
		tree.Insert(exInt(r.Int()))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f := recurse()
		tree.Map(PostOrder, f)

	}

}
