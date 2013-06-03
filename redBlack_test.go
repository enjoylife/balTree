package gotree

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func init() {
}

/* Helpers for tree traversal and testing tree properties */
func printRBNode(n *RBNode) {
	//x := n.Elem
	//fmt.Println("Elem:", x)
}

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

// Max, Min Size tests
func inc(t *testing.T) func(n *RBNode) {
	var prior int = -1
	return func(n *RBNode) {
		if prior < int(n.Elem.(exInt)) {
			//fmt.Println("VALUE: ", value.(int))
			prior = int(n.Elem.(exInt))
		} else {
			t.Errorf("Prior: %d, Current: %d", prior, n.Elem)
		}
	}
}
func TestMaxInsert(t *testing.T) {
	t.Parallel()
	tree := &RBTree{}

	if tree.Max() != nil {
		fmt.Println(tree.Max())
		t.Errorf("Max not working")
	}
	for i := 0; i < iters; i++ {
		tree.Insert(exInt(i))
		if tree.Max() != exInt(i) {
			t.Errorf("Max not updateing")
		}
	}

}
func TestMinInsert(t *testing.T) {
	t.Parallel()
	tree := &RBTree{}

	if tree.Min() != nil {
		fmt.Println(tree.Min())
		t.Errorf("Min not working")
	}
	for i := iters; i > 0; i-- {
		tree.Insert(exInt(i))
		if tree.Min() != exInt(i) {
			t.Errorf("Min not updateing")
		}
	}

}
func TestSizeInsert(t *testing.T) {
	t.Parallel()

	tree := &RBTree{}
	for i := 0; i < iters; i++ {
		tree.Insert(exInt(i))
		if tree.Size != i+1 {
			t.Errorf("Size not correctly updateing")
		}
	}

}

// Edge case tests for insert
func TestBasicInsert(t *testing.T) {

	t.Parallel()
	var old Interface
	items := []exStruct{exStruct{0, "0"},
		exStruct{2, "2"}, exStruct{2, "3"}}
	tree := &RBTree{}

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
	if tree.Min() != items[0] {
		fmt.Println(tree.Min())
		t.Errorf("Min not working")
	}
	if tree.Max() != items[1] {
		fmt.Println(tree.Max())
		t.Errorf("Max not working")
	}
}
func TestRandomInsert(t *testing.T) {

	t.Parallel()
	r := rand.New(rand.NewSource(int64(5)))
	tree := &RBTree{}

	tree.Map(PostOrder, printRBNode)
	for i := 0; i < iters; i++ {
		item := r.Int()
		tree.Insert(exInt(item))
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}
	tree.Map(InOrder, inc(t))
}

func TestSearch(t *testing.T) {

	var elem Interface
	t.Parallel()
	tree := &RBTree{}

	elem = tree.Search(exInt(1))
	if elem != nil {
		t.Errorf("Not minding empty tree")
	}

	for i := 0; i < iters; i++ {
		tree.Insert(exInt(i))
	}
	elem = tree.Search(exInt(1))
	if elem == nil {
		t.Errorf("Not minding nil key's")
	}

	tree.Map(InOrder, inc(t))
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

// remove, min and max tests
func TestRemove(t *testing.T) {

	t.Parallel()
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

	max := 100
	for i := 0; i < max; i++ {
		tree.Insert(exStruct{i, strconv.Itoa(i)})
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
func TestMaxRemove(t *testing.T) {
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
func TestMinRemove(t *testing.T) {
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
func TestSizeRemove(t *testing.T) {
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
func TestRandomRemove(t *testing.T) {

	t.Parallel()
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
func TestIterIn(t *testing.T) {

	t.Parallel()
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
	var check *RBNode
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
	count = 0
	for i, n := 0, tree.IterInit(PreOrder); n != nil; i, n = i+1, tree.Next() {

		count++
		if preOrder[i] != string(n.Elem.(exString)) {
			t.Errorf("Elems are in wrong order Got:%s, Exp: %s", n.Elem, preOrder[i])
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
		if postOrder[i] != string(n.Elem.(exString)) {
			t.Errorf("Values are in wrong order Got:%s, Exp: %s", n.Elem, postOrder[i])
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
func TestTraversal(t *testing.T) {
	t.Parallel()

	tree := &RBTree{}
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for _, v := range items {
		tree.Insert(exString(v))
	}
	if !isBalanced(tree) {
		t.Errorf("Tree is not balanced")
	}
	//tree.Map(InOrder, printRBNode)
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

func BenchmarkMapInsert(b *testing.B) {

	b.StopTimer()
	//r := rand.New(rand.NewSource(int64(5)))
	//m := make(map[int]int)
	m := make(map[int]exInt)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//a := r.Int()
		//m[a] = a
		m[i] = (exInt(b.N - i))
	}

}
func BenchmarkInsert(b *testing.B) {

	b.StopTimer()
	tree := &RBTree{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(exInt(b.N - i))
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
		tree.Remove(exInt(b.N - i))
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

		for i, n := 0, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {
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

		for i, n := 0, tree.IterInit(PreOrder); n != nil; i, n = i+1, tree.Next() {
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

		for i, n := 0, tree.IterInit(PostOrder); n != nil; i, n = i+1, tree.Next() {
			sum += int(n.Elem.(exInt))
		}
	}

}

func recurse() func(n *RBNode) {
	var sum int = 0
	return func(n *RBNode) {
		sum += int(n.Elem.(exInt))
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
