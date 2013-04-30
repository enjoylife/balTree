package gorbtree

import (
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"gotree"
	"math/rand"
	"testing"
)

var _ = fmt.Printf

func testCmp(a interface{}, b interface{}) gotree.Direction {
	switch result := (a.(int) - b.(int)); {
	case result > 0:
		return gotree.GT
	case result < 0:
		return gotree.LT
	case result == 0:
		return gotree.EQ
	default:
		panic("Invalid Compare function Result")
	}
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
func printNode(key interface{}, value interface{}) {
	fmt.Println("VALUE: ", value.(int))
}

func TestInsert(t *testing.T) {

	r := rand.New(rand.NewSource(int64(5)))
	tree := New(testCmp)
	iters := 17
	for i := 0; i < iters; i++ {
		item := r.Int()
		item = i
		//fmt.Printf("Iteration %d\n", i)
		tree.Insert(item, item)
	}
	tree.Traverse(printNode)
	//order := tree.TraverseTest(inc(t))
	//fmt.Printf("Len: %d \n", len(order))
	/*for i := 0; i < len(order); i++ {
		if order[i] != nil {
			fmt.Println(order[i].value.(int))
			//a := order[i].value.(int)
		}
	}*/
}

func (t *RbTree) TraverseTest(f gotree.IterFunc) []*rbNode {
	node := t.root
	order := []*rbNode{nil}
	stack := []*rbNode{nil}
	for len(stack) != 1 || node != nil {
		if node != nil {
			stack = append(stack, node)
			order = append(order, node)
			node = node.right
		} else {
			stackIndex := len(stack) - 1
			node = stack[stackIndex]
			f(node.key, node.value)
			stack = stack[0:stackIndex]
			//fmt.Println("stack size", len(stack))
			node = node.left
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
	//spew.Dump(tree.root)
	fmt.Println("Height", tree.Height)

}
