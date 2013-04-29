package gorbtree

import (
	"../"
	"fmt"
	"math/rand"
	"testing"
)

var _ = fmt.Printf

func testCmp(a interface{}, b interface{}) int {
	return a.(int) - b.(int)
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
	tree := New(testCmp)
	iters := 10
	for i := 0; i < iters; i++ {
		item := r.Int()
		item = i
		//fmt.Printf("Iteration %d\n", i)
		tree.Insert(item, item)
	}
	order := tree.TraverseTest(inc(t))
	fmt.Printf("Len: %d \n", len(order))
	for i := 0; i < len(order); i++ {
		if order[i] != nil {
			//a := order[i].value.(int)
		}
	}
}

func (t *RbTree) TraverseTest(f gotree.IterFunc) []*rbNode {
	node := t.root
	order := []*rbNode{nil}
	stack := []*rbNode{nil}
	for len(stack) != 1 || node != nil {
		if node != nil {
			stack = append(stack, node)
			order = append(order, node)
			node = node.left
		} else {
			stackIndex := len(stack) - 1
			node = stack[stackIndex]
			f(node.key, node.value)
			stack = stack[0:stackIndex]
			//fmt.Println("stack size", len(stack))
			node = node.right
		}
	}
	return order
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
	tree := New(testCmp)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		old, ok := tree.Insert(r.Int(), r.Int())
		if !ok {
			_ = old
		}
	}

}
